package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"io/ioutil"
	"path/filepath"

	"github.com/xhaklaaa/go-telegrambot-images/internal/image_creator"
	"github.com/xhaklaaa/go-telegrambot-images/pkg/parser"
	"github.com/xhaklaaa/go-telegrambot-images/pkg/telegram"
	"gopkg.in/yaml.v2"
)

const (
	quotesFile  = "internal/image_creator/quotes.json"
	flagFile    = "parser_completed.flag"
	configFile  = "configs/configs.yaml"
	resultedDir = "Resulted"
)

type Config struct {
	ChannelID int64  `yaml:"channelID"`
	BotToken  string `yaml:"botToken"`
	AccessKey string `yaml:"accessKey"`
}

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = configFile
	}

	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	bot, err := telegram.NewBot(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	if parser.ShouldRunParser(flagFile) {
		err := parser.ParseQuotes()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Quotes saved to quotes.json")

		parser.CreateFlagFile(flagFile)
	} else {
		fmt.Println("Parser has already been run. Skipping parsing.")
	}

	quotes, err := parser.ReadQuotes(quotesFile)
	if err != nil {
		log.Fatal(err)
	}

	accessKey := config.AccessKey
	sendQuotesToTelegram(bot, config.ChannelID, quotes, accessKey)
}

func readConfig(filename string) (*Config, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func sendQuotesToTelegram(bot *telegram.Bot, channelID int64, quotes []parser.Quote, accessKey string) {
	lastTags := make(map[string]bool)

	if err := os.MkdirAll(resultedDir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %s", err)
	}

	for len(quotes) > 0 {
		quote, err := parser.GetUniqueQuote(quotes, lastTags)
		if err != nil {
			log.Printf("Failed to get unique quote: %s", err)
			continue
		}

		// Генерация изображения для цитаты
		imagePath := filepath.Join(resultedDir, fmt.Sprintf("output_image_%s.jpg", quote.Author))
		err = image_creator.GenerateImageForQuote(image_creator.Quote{
			Author: quote.Author,
			Quote:  quote.Text,
			Tags:   quote.Tags,
		}, imagePath, accessKey)
		if err != nil {
			log.Printf("Failed to generate image for quote: %s", err)
		} else {
			log.Printf("Generated image for quote: %s", imagePath)
			err = bot.SendPhotoWithTags(channelID, imagePath, quote.Tags)
			if err != nil {
				log.Printf("Failed to send photo: %s", err)
			}
		}

		for _, tag := range quote.Tags {
			lastTags[tag] = true
		}
		quotes = parser.RemoveQuote(quotes, quote)

		err = parser.WriteQuotes(quotesFile, quotes)
		if err != nil {
			log.Printf("Failed to write quotes to file: %s", err)
		}

		time.Sleep(30 * time.Second)
	}
}
