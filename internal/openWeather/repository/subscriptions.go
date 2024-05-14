package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/internal/logger"
	"main.go/internal/openWeather/model"
)

type SLogger interface {
	SLog() *slog.Logger
}

type SLog struct{}

func (L SLog) SLog() *slog.Logger { return logger.GetLogger() }

var sLogger = new(SLog).SLog()

type SubscriptionRepo struct {
	MongoCollection *mongo.Collection
}

func (r *SubscriptionRepo) AssigneeSubscription(sub *model.Subscriptions) (interface{}, error) {
	res, err := r.MongoCollection.InsertOne(context.Background(), sub)
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	return res.InsertedID, nil
}

func (r *SubscriptionRepo) GetAllAssignedSubscriptions() ([]model.Subscriptions, error) {
	res, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	var sub []model.Subscriptions
	err = res.All(context.Background(), &sub)
	if err != nil {
		sLogger.Error("Send message error", "err", err)
	}

	return sub, nil
}

func (r *SubscriptionRepo) DeleteSubscription(id interface{}) (int64, error) {
	res, err := r.MongoCollection.DeleteOne(context.Background(), bson.D{{Key: "user_id", Value: id}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
