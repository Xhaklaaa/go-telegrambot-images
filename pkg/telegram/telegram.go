package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{bot: bot}, nil
}

func (b *Bot) SendPhotoWithTags(channelID int64, filePath string, tags []string) error {
	photo := tgbotapi.NewPhotoUpload(channelID, filePath)
	caption := ""
	for _, tag := range tags {
		caption += "#" + tag + " "
	}
	photo.Caption = caption

	_, err := b.bot.Send(photo)
	return err
}
