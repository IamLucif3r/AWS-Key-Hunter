package pkg

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func SearchGithub(githubToken, sort string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	query := `AKIA filename:.env OR filename:.ini OR filename:.yml OR filename:.yaml OR filename:.json`
	opts := &github.SearchOptions{
		Sort:      sort,
		Order:     "desc",
		TextMatch: true,
	}

	if shouldBackoff(ctx, client) {
		return
	}

	log.Printf("🔍 Searching GitHub for AWS keys (sort: %s)...", strings.ToUpper(sort))

	results, _, err := client.Search.Code(ctx, query, opts)
	if err != nil {
		log.Printf("❌ GitHub search error: %v", err)
		return
	}

	if len(results.CodeResults) == 0 {
		log.Println("✅ No matching files found in this cycle.")
		return
	}

	log.Printf("📦 Found %d possible files. Scanning for credentials...", len(results.CodeResults))

	for _, file := range results.CodeResults {
		checkFileContent(ctx, client, &file)
	}
}

func shouldBackoff(ctx context.Context, client *github.Client) bool {
	rateLimits, _, err := client.RateLimits(ctx)
	if err != nil {
		log.Printf("⚠️  Could not fetch rate limits: %v", err)
		return false
	}

	remaining := rateLimits.Core.Remaining
	resetTime := rateLimits.Core.Reset.Time
	resetIn := time.Until(resetTime)

	if remaining < 5 {
		log.Printf("⏳ GitHub API rate limit nearly exhausted. Remaining: %d. Waiting %v until reset (%s)...",
			remaining, resetIn.Truncate(time.Second), resetTime.Format(time.RFC1123))
		time.Sleep(resetIn + 5*time.Second)
		return true
	}

	return false
}
