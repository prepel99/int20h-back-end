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
}

type UserStore struct {
	DB *mongo.Client
}

type WorkOut struct {
	Type    string        `json:"type" bson:"type"`
	Results []sensorValue `json:"results" bson:"results"`
}

type sensorValue struct {
	Value int   `json:"value" bson:"value"`
	Date  int64 `json:"date" bson:"date"`
}

type User struct {
	ID         string             `json:"ID" bson:"_id"`
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

	objID, _ := primitive.ObjectIDFromHex(id)
	fmt.Println(objID)
	filter := bson.D{{"_id", objID}}

	user := User{}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserStore) SaveOneExercise(userID string, exercise WorkOut) (User, error) {
	collection := u.DB.Database("sensors").Collection("users")

	objID, _ := primitive.ObjectIDFromHex(userID)
	fmt.Println(objID)
	filter := bson.D{{"_id", objID}}

	user := User{}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
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

func (u *UserStore) GetAllUsers() ([]User, error) {
	collection := u.DB.Database("sensors").Collection("users")

	fmt.Println(111)
	options := options.Find()
	filter := bson.M{}

	// Here's an array in which you can store the decoded documents
	var results []User

	// Passing nil as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
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
