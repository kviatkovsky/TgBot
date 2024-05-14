package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscriptions struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	UserID                interface{}        `json:"user_id" bson:"user_id"`
	WeatherSubscriptionID interface{}        `json:"weather_subscription_id" bson:"weather_subscription_id"`
}
