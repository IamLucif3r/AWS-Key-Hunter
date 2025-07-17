# AWS-Key-Hunter

AWS Key Hunter is a powerful and automated tool that scans GitHub repositories for exposed AWS keys. It continuously monitors commits, detects AWS secrets in both base64 and plaintext formats, and alerts users about potential security risks on Discord.

## Features 🚀

- **Real-time Monitoring:** Continuously tracks new commits in GitHub repositories.
- **Comprehensive AWS Key Detection:** Identifies AWS credentials in both plaintext and base64-encoded formats.
- **Automated Scanning:** Performs scheduled scans to uncover exposed AWS keys.
- **Resource-Efficient & Secure:** Designed for low resource consumption and deployed securely via Docker.
- **Discord Integration:** Instantly notifies users of detected credentials through Discord alerts.

## Installation 📥

Create a `.env` file and add your **Github** token and your **Discord** Server's web hook in the file. 

### Using Docker

Build the Docker image
```bash
docker build -t aws-key-scanner .
```
Run the container
```bash
docker run --rm -d --name aws-scanner aws-key-scanner

# Recommended: Don't bulid your dockerfile with .env file into it, instead pass them as environment variables:
docker run -e GITHUB_TOKEN="<your-github-token>" -e DISCORD_WEBHOOK="<discord-webhook-url>" aws-key-hunter:latest
```

## Usage 🛠

Running Locally
```bash
echo "GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" > .env
echo "DISCORD_WEBHOOK=https://discord.com/api/webhooks/..." >> .env

go run main.go
```

## Disclaimer ⚠️

This tool was created for educational and experimental purposes only. They are not intended to be used for malicious activities or to harm others in any way. I do not endorse or encourage the use of this tool or information for illegal, unethical, or harmful actions.

By using this tool, you agree to accept full responsibility for any consequences that may arise from its use. I will not be held accountable for any damages, losses, or legal repercussions resulting from the misuse of this tool or the information provided.

Use at your own risk.

## Contributing 🤝

Contributions are welcome! Feel free to open an issue or submit a PR.