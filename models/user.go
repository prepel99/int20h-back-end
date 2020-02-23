package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserStorer interface {
	RegisterUser(user User) (string, error)
	GetUser(id string) (User, error)
	SaveOneExercise(userID string, exercise WorkOut) (User, error)
	GetAllUsers() ([]User, error)
	SaveExerciseWithSensor(userID string, exercise WorkOut, secondsAgo int) (User, error)
}

type UserStore struct {
	DB *mongo.Client
}

type WorkOut struct {
	Type    string        `json:"type" bson:"type"`
	Results []SensorValue `json:"results" bson:"results"`
}

type SensorValue struct {
	AvgSensorValue int   `json:"sensorValue" bson:"sensorValue"`
	Value          int   `json:"value" bson:"value"`
	Date           int64 `json:"date" bson:"date"`
}

type User struct {
	Name       string             `json:"name" bson:"name"`
	Workouts   map[string]WorkOut `json:"workouts" bson:"workouts"`
	Updated_at int64              `json:"updated_at" bson:"updated_at"`
}

func (u *UserStore) RegisterUser(user User) (string, error) {
	collection := u.DB.Database("sensors").Collection("users")

	now := time.Now().Unix()
	user.Updated_at = now

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}
	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("cannot parse ObjectID")
	}
	return newID.Hex(), nil
}

func (u *UserStore) GetUser(id string) (User, error) {
	collection := u.DB.Database("sensors").Collection("users")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User{}, err

	}
	filter := bson.D{{"_id", objID}}

	user := User{}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserStore) SaveOneExercise(userID string, exercise WorkOut) (User, error) {
	collection := u.DB.Database("sensors").Collection("users")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return User{}, err

	}
	filter := bson.D{{"_id", objID}}

	user := User{}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	userResults := user.Workouts[exercise.Type].Results
	for _, v := range exercise.Results {
		if v.Date == 0 {
			now := time.Now().Unix()
			v.Date = now
		}
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

func (u *UserStore) SaveExerciseWithSensor(userID string, exercise WorkOut, secondsAgo int) (User, error) {
	collection := u.DB.Database("sensors").Collection("users")
	collectionSensors := u.DB.Database("sensors").Collection("SensorsData")

	objID, err := primitive.ObjectIDFromHex("5e51966d09eaf8c6d663ff3c")
	if err != nil {
		return User{}, err
	}
	filter := bson.D{{"_id", objID}}

	user := User{}
	err = collectionSensors.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}

	exerciseValues := user.Workouts[exercise.Type]
	accExercises := make([]SensorValue, 1)

	now := time.Now()
	backTime := time.Second * time.Duration(-secondsAgo)
	past := now.Add(backTime).Unix()
	avgValue := 0
	for _, v := range exerciseValues.Results {
		if v.Date > past {
			accExercises = append(accExercises, v)
			avgValue += v.Value
		}
	}
	avgValue = avgValue / len(accExercises)

	objID, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		return User{}, err
	}

	filter = bson.D{{"_id", objID}}

	user = User{}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	fmt.Println(1111)

	userResults := user.Workouts[exercise.Type].Results
	for _, v := range exercise.Results {
		if v.Date == 0 {
			now := time.Now().Unix()
			v.Date = now
		}
		v.AvgSensorValue = avgValue
		userResults = append(userResults, v)
	}
	fmt.Println(1111)

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

func (u *UserStore) GetAllUsers() ([]User, error) {
	collection := u.DB.Database("sensors").Collection("users")

	options := options.Find()
	filter := bson.M{}

	var results []User

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		fmt.Println(elem)
		results = append(results, elem)
	}
	return results, nil
}
