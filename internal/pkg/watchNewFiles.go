package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/go-github/github"
)

func checkFileContent(ctx context.Context, client *github.Client, file *github.CodeResult) {
	content, err := fetchFileContent(ctx, client, file)
	if err != nil {
		log.Printf("❌ Failed to fetch content: %v", err)
		return
	}

	credsList := extractAWSKeys(content)
	if len(credsList) == 0 {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, creds := range credsList {
		accessKey := creds["access_key"]
		secretKey := creds["secret_key"]

		if !looksLikeAWSKey(accessKey, secretKey) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(accessKey, secretKey string) {
			defer wg.Done()
			defer func() { <-sem }()

			if validateAWSKeys(accessKey, secretKey) {
				repo := file.Repository.GetFullName()
				url := file.GetHTMLURL()
				log.Printf("🚨 Valid AWS Key Detected | Repo: %s | File: %s", repo, url)
				sendDiscordAlert(repo, url, []string{maskKey(accessKey)})
			}
		}(accessKey, secretKey)
	}

	wg.Wait()
}

func fetchFileContent(ctx context.Context, client *github.Client, file *github.CodeResult) (string, error) {
	repo := file.GetRepository()
	owner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	filePath := file.GetPath()

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repoName, filePath, nil)
	if err != nil {
		return "", fmt.Errorf("GitHub API error fetching file content: %v", err)
	}

	if fileContent == nil {
		return "", errors.New("file content is nil")
	}

	contentStr, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("error retrieving file content: %v", err)
	}

	return strings.TrimSpace(contentStr), nil
}

func extractAWSKeys(content string) []map[string]string {
	awsKeys := []map[string]string{}

	re := regexp.MustCompile(`(?i)(aws_access_key_id.*?(AKIA[0-9A-Z]{16})).*?(aws_secret_access_key.*?([a-zA-Z0-9/+]{40}))`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 5 {
			awsKeys = append(awsKeys, map[string]string{
				"access_key": match[2],
				"secret_key": match[4],
			})
		}
	}

	if len(awsKeys) == 0 {
		accessKeyPattern := regexp.MustCompile(`AKIA[0-9A-Z]{16}`)
		secretKeyPattern := regexp.MustCompile(`[a-zA-Z0-9/+]{40}`)

		accessKeys := accessKeyPattern.FindAllString(content, -1)
		secretKeys := secretKeyPattern.FindAllString(content, -1)

		for i := range accessKeys {
			if i < len(secretKeys) {
				awsKeys = append(awsKeys, map[string]string{
					"access_key": accessKeys[i],
					"secret_key": secretKeys[i],
				})
			}
		}
	}

	return awsKeys
}

func looksLikeAWSKey(accessKey, secretKey string) bool {
	return len(accessKey) == 20 && len(secretKey) == 40
}

func validateAWSKeys(accessKey, secretKey string) bool {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return false
	}

	stsClient := sts.NewFromConfig(cfg)
	_, err = stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	return err == nil
}

func maskKey(key string) string {
	if len(key) < 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

func sendDiscordAlert(repo, url string, keys []string) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	if webhookURL == "" {
		log.Println("⚠️  DISCORD_WEBHOOK not set, skipping alert.")
		return
	}

	message := map[string]string{
		"content": fmt.Sprintf("🚨 *AWS Key Leak Detected!*\n**Repo**: `%s`\n**URL**: %s\n**Keys**: `%v`", repo, url, keys),
	}
	jsonData, _ := json.Marshal(message)

	req, _ := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to send alert to Discord: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Println("📣 Alert sent to Discord.")
}
