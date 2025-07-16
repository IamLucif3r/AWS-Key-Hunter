package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/iamlucif3r/aws-key-hunter/internal/pkg"
)

const (
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  No .env file found, falling back to environment variables...")
	}
}

func main() {
	fmt.Println(Red + `
┏┓┓ ┏┏┓  ┓┏┓      ┓┏     		
┣┫┃┃┃┗┓━━┃┫ ┏┓┓┏━━┣┫┓┏┏┓╋┏┓┏┓	
┛┗┗┻┛┗┛  ┛┗┛┗ ┗┫  ┛┗┗┻┛┗┗┗ ┛  	
               ┛   v1.0.1      	
` + Reset)

	log.Println(Yellow + "🚀 Starting AWS Key Hunter..." + Reset)

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("❌ GITHUB_TOKEN is not set in environment")
	}

	go func() {
		for {
			pkg.SearchGithub(githubToken, "updated")
			time.Sleep(2 * time.Minute)
		}
	}()

	for {
		pkg.SearchGithub(githubToken, "indexed")
		time.Sleep(5 * time.Minute)
	}
}
