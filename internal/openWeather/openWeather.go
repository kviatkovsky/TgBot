package openWeather

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/internal/logger"
)

var DB *mongo.Client = InitMongoDB()

type SLogger interface {
	SLog() *slog.Logger
}

type SLog struct{}

type WeatherResponse struct {
	Coord      *Coord        `json:"coord"`
	Weather    []WeatherInfo `json:"weather"`
	Base       string        `json:"base"`
	Main       *MainInfo     `json:"main"`
	Visibility int           `json:"visibility"`
	Wind       *WindInfo     `json:"wind"`
	Clouds     *CloudInfo    `json:"clouds"`
	Dt         int64         `json:"dt"`
	Sys        *SysInfo      `json:"sys"`
	Timezone   int           `json:"timezone"`
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Cod        int           `json:"cod"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type WeatherInfo struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type MainInfo struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type WindInfo struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type CloudInfo struct {
	All int `json:"all"`
}

type SysInfo struct {
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

func (r *factRepository) GetWeatherFromApi() WeatherResponse {
	var response WeatherResponse
	req, err := http.NewRequest(http.MethodGet, r.address, nil)
	if err != nil {
		sLogger.Error(err.Error())
	}

	res, err := r.client.Do(req)
	if err != nil {
		sLogger.Error(err.Error())
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		sLogger.Error(err.Error())
	}

	return response
}

func (L SLog) SLog() *slog.Logger { return logger.GetLogger() }

func GetWeather(lat float64, lon float64) WeatherResponse {
	url := GetQueryUrl(lat, lon)
	repo := newFactRepository(url)

	return repo.GetWeatherFromApi()
}

type factRepository struct {
	address string
	client  *http.Client
}

func newFactRepository(addr string) *factRepository {
	return &factRepository{
		address: addr,
		client:  http.DefaultClient,
	}
}

func InitMongoDB() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(GetConfigs().MongoUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		sLogger.Error(err.Error())
	}

	if err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		sLogger.Error(err.Error())
	}

	return client
}

func GetUserCollection() *mongo.Collection {
	return DB.Database(GetConfigs().DBName).Collection(GetConfigs().UserCollectionName)
}

func GetWeatherSubscriptionCollection() *mongo.Collection {
	return DB.Database(GetConfigs().DBName).Collection(GetConfigs().WeatherSubscriptionCollectionName)
}

func GetSubscriptionCollection() *mongo.Collection {
	return DB.Database(GetConfigs().DBName).Collection(GetConfigs().SubscriptionCollectionName)
}
