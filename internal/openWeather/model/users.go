package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	ChatId int64              `json:"chat_id" bson:"chat_id"`
	Lat    float64            `json:"lat" bson:"lat"`
	Lon    float64            `json:"lon" bson:"lon"`
}
