package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SensorStorer interface {
	RegisterUser(user User) (string, error)
	SaveOneExercise(userID string, exercise WorkOut) (User, error)
}

type SensorStore struct {
	DB *mongo.Client
}

func (s *SensorStore) RegisterUser(user User) (string, error) {
	collection := s.DB.Database("sensors").Collection("SensorsData")

	now := time.Now().Unix()
	user.Updated_at = now

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}
	fmt.Println(result)
	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("cannot parse ObjectID")
	}
	return newID.Hex(), nil
}

func (u *SensorStore) SaveOneExercise(userID string, exercise WorkOut) (User, error) {
	collection := u.DB.Database("sensors").Collection("SensorsData")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return User{}, nil
	}
	filter := bson.D{{"_id", objID}}

	user := User{}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	userResults := user.Workouts[exercise.Type].Results
	for _, v := range exercise.Results {
		now := time.Now().Unix()
		v.Date = now
		userResults = append(userResults, v)
	}

	workout := user.Workouts[exercise.Type]
	workout.Results = userResults

	user.Workouts[exercise.Type] = workout

	doc, err := toDoc(user)
	if err != nil {
		return User{}, err
	}
	update := bson.D{{"$set", *doc}}
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return User{}, nil
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	return user, nil
}
