package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/internal/openWeather/model"
)

type WeatherSubscriptionRepo struct {
	MongoCollection *mongo.Collection
}

func (r *WeatherSubscriptionRepo) InsertWeatherSubscription(wthSub *model.WeatherSubscription) (interface{}, error) {
	res, err := r.MongoCollection.InsertOne(context.Background(), wthSub)
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	return res.InsertedID, nil
}

func (r *WeatherSubscriptionRepo) GetAllWeatherSubscriptions() ([]model.WeatherSubscription, error) {
	res, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	var wthSub []model.WeatherSubscription
	err = res.All(context.Background(), &wthSub)
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	return wthSub, nil
}

func (r *WeatherSubscriptionRepo) GetWeatherSubscriptionsByName(weatherSubName string) ([]model.WeatherSubscription, error) {
	res, err := r.MongoCollection.Find(context.Background(), bson.D{{Key: "name", Value: weatherSubName}})
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	var wthSub []model.WeatherSubscription
	err = res.All(context.Background(), &wthSub)
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	return wthSub, nil
}

func (r *WeatherSubscriptionRepo) GetWeatherSubscriptionById(id interface{}) ([]model.WeatherSubscription, error) {
	res, err := r.MongoCollection.Find(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	var wthSub []model.WeatherSubscription
	err = res.All(context.Background(), &wthSub)
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	return wthSub, nil
}
