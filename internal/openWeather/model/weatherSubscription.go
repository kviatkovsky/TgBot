package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type WeatherSubscription struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	NotificationTime string             `json:"notification_time" bson:"notification_time"`
}
