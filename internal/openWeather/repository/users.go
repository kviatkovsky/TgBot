package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/internal/openWeather/model"
)

type UsersRepo struct {
	MongoCollection *mongo.Collection
}

func (r *UsersRepo) InsertUser(usr *model.Users) (interface{}, error) {
	res, err := r.MongoCollection.InsertOne(context.Background(), usr)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (r *UsersRepo) FindUserByChatAndName(chatID int64, userName string) (*model.Users, error) {
	var usr model.Users

	err := r.MongoCollection.FindOne(context.Background(), bson.D{
		{Key: "chat_id", Value: chatID},
		{Key: "name", Value: userName},
	}).Decode(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *UsersRepo) GetUserById(id interface{}) (*model.Users, error) {
	var usr model.Users

	err := r.MongoCollection.FindOne(context.Background(), bson.D{
		{Key: "_id", Value: id},
	}).Decode(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *UsersRepo) DeleteUser(userID string) (int64, error) {
	res, err := r.MongoCollection.DeleteOne(context.Background(), bson.D{{Key: "user_id", Value: userID}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
