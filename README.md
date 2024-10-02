# Telegrambot-images
Telegrambot-images is a Go-based project that automates the process of generating and sharing inspirational quotes with images in a Telegram channel. 
The project consists of three main steps: parsing quotes, generating images with quotes, and sending the final images to a Telegram channel.
## Requirements
•	Go 1.20 or higher

• [Go Graphics](https://github.com/fogleman/gg)

• [Goquery](https://github.com/PuerkitoBio/goquery)

• [Telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)

## Features

**Quote Parsing:** Extracts quotes from a predefined source and saves them in a JSON file.

**Image Generation:** Fetches a random image from Unsplash API, overlays the quote and author's name, and saves the final image.

**Telegram Integration:** Sends the generated image to a specified Telegram channel.

## Prerequisites

Before you begin, ensure you have met the following requirements:

**Go:** Install Go from [golang.org](golang.org).

**Telegram Bot:** Create a Telegram bot and obtain the bot token.

**Unsplash API:** Sign up for an Unsplash developer account and obtain an API key.

**Environment Variables:** Set up environment variables for configuration.

## Installation
• **Clone the repository:**
```
git clone https://github.com/xhaklaaa/go-telegrambot-images.git
cd go-telegrambot-images
```

• **Install dependencies:**
```
go mod download
```
## Configuration
• **Environment Variables: **Create a configs/configs.yaml file in the root directory with the following content:
```
channelID: your_telegram_channel_id
botToken: your_telegram_bot_token
accessKey: your_unsplash_api_key
```
## Usage
Run the project
```
go run cmd/main/main.go
```
