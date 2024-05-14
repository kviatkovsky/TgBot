package bot

import (
	tgbotapi "github.com/crocone/telegram-bot-api"
	"main.go/internal/configs"
)

var Bot *tgbotapi.BotAPI

func InitBot() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(configs.GetConfigs().TelegramToken)
	if err != nil {
		sLogger.Error("Failed to initialize Telegram bot API", "err", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := Bot.GetUpdatesChan(u)

	if err != nil {
		sLogger.Error("Failed to get updates", "err", err)
	}

	HandleWeatherSubscription()
	HandleUpdates(updates)
}
