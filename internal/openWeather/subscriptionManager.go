package openWeather

import (
	"fmt"
	"time"

	"main.go/internal/openWeather/model"
	"main.go/internal/openWeather/repository"
)

var userCollection = GetUserCollection()
var weatherSubscriptionCollection = GetWeatherSubscriptionCollection()
var subscriptionsCollection = GetSubscriptionCollection()

func CreateUserIfNotExist(user model.Users) interface{} {
	userRepo := repository.UsersRepo{MongoCollection: userCollection}

	existingUser, _ := userRepo.FindUserByChatAndName(user.ChatId, user.Name)

	if existingUser == nil {
		insertUser, err := userRepo.InsertUser(&user)
		if err != nil {
			sLogger.Error(err.Error())
		}

		return insertUser
	}

	return existingUser.ID.Hex()
}

func UnSubscribe(chatID int64, userName string) int64 {
	subRepo := repository.SubscriptionRepo{MongoCollection: subscriptionsCollection}
	userRepo := repository.UsersRepo{MongoCollection: userCollection}

	existingUser, _ := userRepo.FindUserByChatAndName(chatID, userName)

	subscription, err := subRepo.DeleteSubscription(existingUser.ID)
	if err != nil {
		sLogger.Error(err.Error())
	}

	return subscription
}

func Subscribe(user model.Users) {
	weatherSubscriptionID := CreateWeatherSubscription() // Temporary solution. TODO Move subscription creation from this place
	userId := CreateUserIfNotExist(user)
	if weatherSubscriptionID != nil || userId != nil {
		AssigneeUserToSubscription(userId, weatherSubscriptionID)
	}
}

func AssigneeUserToSubscription(userId interface{}, weatherSubscriptionID interface{}) {
	subRepo := repository.SubscriptionRepo{MongoCollection: subscriptionsCollection}

	subs := model.Subscriptions{
		WeatherSubscriptionID: weatherSubscriptionID,
		UserID:                userId,
	}

	_, err := subRepo.AssigneeSubscription(&subs)
	if err != nil {
		sLogger.Error(err.Error())
	}
}

func CreateWeatherSubscription() interface{} {
	wSubRepo := repository.WeatherSubscriptionRepo{MongoCollection: weatherSubscriptionCollection}

	existingWeatherSubscription := GetSubscriptionIfExist(wSubRepo)
	if existingWeatherSubscription != false {
		return existingWeatherSubscription
	}

	timeLayout := "15:04:05 -0700 MST"

	inputTimeStr := "15:00:00 +0200 CEST"

	inputTime, err := time.Parse(timeLayout, inputTimeStr)
	if err != nil {
		sLogger.Error(err.Error())
	}

	currentDate := time.Now().Format("2006-01-02")

	notificationTime := fmt.Sprintf("%s %s", currentDate, inputTime.Format("15:04:05"))
	wthSubModel := model.WeatherSubscription{
		NotificationTime: notificationTime,
		Name:             "WeatherSubscription",
	}

	weatherSubscriptionID, err := wSubRepo.InsertWeatherSubscription(&wthSubModel)
	if err != nil {
		sLogger.Error(err.Error())
	}

	return weatherSubscriptionID
}

func GetSubscriptionIfExist(wSubRepo repository.WeatherSubscriptionRepo) interface{} {
	weatherSubscriptions, err := wSubRepo.GetAllWeatherSubscriptions()
	if err != nil {
		sLogger.Error(err.Error())
		return false
	}

	if len(weatherSubscriptions) == 0 {
		return false
	}

	return weatherSubscriptions[0].ID.Hex()
}

func GetAllSubscriptions() ([]model.Subscriptions, error) {
	subRepo := repository.SubscriptionRepo{MongoCollection: subscriptionsCollection}

	return subRepo.GetAllAssignedSubscriptions()
}

func GetUserById(id interface{}) (*model.Users, error) {
	userRepo := repository.UsersRepo{MongoCollection: userCollection}

	return userRepo.GetUserById(id)
}

func GetWeatherSubscriptionById(id interface{}) ([]model.WeatherSubscription, error) {
	wSubRepo := repository.WeatherSubscriptionRepo{MongoCollection: weatherSubscriptionCollection}

	return wSubRepo.GetWeatherSubscriptionById(id)
}
