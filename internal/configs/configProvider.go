package configs

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"main.go/internal/logger"
)

type Config struct {
	TelegramToken      string
	HolidayBotApiKey   string
	HolidayBotApiEmail string
}

type SLog struct{}

type SLogger interface {
	SLog() *slog.Logger
}

func (L SLog) SLog() *slog.Logger { return logger.GetLogger() }

var sLogger = new(SLog).SLog()

func GetConfigs() Config {
	cfg := Config{}
	err := godotenv.Load(".env")
	if err != nil {
		sLogger.Error(".env not loaded")
	}

	yamlCfg := readYamlConfigs()
	cfg.TelegramToken = os.Getenv("TG_API_BOT_TOKEN")
	cfg.HolidayBotApiKey = os.Getenv("HOLIDAY_BOT_API_KEY")
	cfg.HolidayBotApiEmail = yamlCfg["holiday_bot_api_email"]

	return cfg
}

func readYamlConfigs() map[string]string {
	cfg := make(map[string]string)

	yamlFile, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		sLogger.Error("main.yml Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		sLogger.Error("Unmarshal: %v", err)
	}

	return cfg
}
