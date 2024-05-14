package openWeather

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	OpenWeatherApiUrl                 string
	OpenWeatherApiKey                 string
	MongoUri                          string
	DBName                            string
	UserCollectionName                string
	SubscriptionCollectionName        string
	WeatherSubscriptionCollectionName string
}

var sLogger = new(SLog).SLog()

func GetConfigs() Config {
	cfg := Config{}
	err := godotenv.Load(".env")
	if err != nil {
		sLogger.Error(".env not loaded")
	}

	cfg.MongoUri = os.Getenv("MONGO_URI")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.UserCollectionName = os.Getenv("USER_COLLECTION_NAME")
	cfg.SubscriptionCollectionName = os.Getenv("SUBSCRIPTIONS_COLLECTION_NAME")
	cfg.WeatherSubscriptionCollectionName = os.Getenv("WEATHER_SUBSCRIPTION_COLLECTION_NAME")

	yamlCfg := readYamlConfigs()

	cfg.OpenWeatherApiUrl = yamlCfg["open_weather_api_url"]
	cfg.OpenWeatherApiKey = yamlCfg["open_weather_api_key"]

	return cfg
}

func GetQueryUrl(lat float64, lon float64) string {
	return fmt.Sprintf("%s?lat=%.2f&lon=%.2f&appid=%s&units=%s",
		GetConfigs().OpenWeatherApiUrl,
		lat,
		lon,
		GetConfigs().OpenWeatherApiKey,
		"metric",
	)
}

func readYamlConfigs() map[string]string {
	cfg := make(map[string]string)

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		sLogger.Error("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		sLogger.Error("Unmarshal: %v", err)
	}

	return cfg
}
