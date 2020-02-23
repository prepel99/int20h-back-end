package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

type Challenge struct {
	FromID     string `json:"fromID" bson:"fromID"`
	ToID       string `json:"toID" bson:"toID"`
	FromName   string `json:"fromName" bson:"fromName"`
	ToName     string `json:"toName" bson:"toName"`
	Type       string `json:"type" bson:"type"`
	Value      int    `json:"value" bson:"value"`
	Done       bool   `json:"done" bson:"done"`
	Result     bool   `json:"result" bson:"result"`
	Date       int64  `json:"date" bson:"date"`
	Updated_at int64  `json:"updated_at" bson:"updated_at"`
}

type SuggestedChallenge struct {
	Type  string
	Value int
	Date  int64
}

type ChallengeStorer interface {
	CreateChallenge(chal Challenge) (string, error)
	GetAllChallenges() ([]Challenge, error)
	GetAllSuggestedChallenges(userID string) ([]SuggestedChallenge, error)
}

type ChallengeStore struct {
	DB *mongo.Client
}

func (c *ChallengeStore) CreateChallenge(chal Challenge) (string, error) {
	collection := c.DB.Database("sensors").Collection("challenges")
	collectionUsers := c.DB.Database("sensors").Collection("users")

	objID, err := primitive.ObjectIDFromHex(chal.ToID)
	if err != nil {
		return "", err
	}
	filter := bson.D{{"_id", objID}}

	toUser := User{}
	fromUser := User{}

	err = collectionUsers.FindOne(context.TODO(), filter).Decode(&toUser)
	if err != nil {
		return "", err
	}

	objID, err = primitive.ObjectIDFromHex(chal.FromID)
	if err != nil {
		return "", err
	}

	filter = bson.D{{"_id", objID}}

	err = collectionUsers.FindOne(context.TODO(), filter).Decode(&fromUser)
	if err != nil {
		return "", err
	}

	chal.ToName = toUser.Name
	chal.FromName = fromUser.Name

	now := time.Now().Unix()
	chal.Date = now
	chal.Updated_at = now

	result, err := collection.InsertOne(context.TODO(), chal)
	if err != nil {
		return "", err
	}

	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("cannot parse ObjectID")
	}

	return newID.Hex(), nil
}

func (c *ChallengeStore) GetAllSuggestedChallenges(userID string) ([]SuggestedChallenge, error) {
	collection := c.DB.Database("sensors").Collection("users")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return []SuggestedChallenge{}, err
	}
	filter := bson.D{{"_id", objID}}

	user := User{}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return []SuggestedChallenge{}, err
	}

	suggestedChlng := make([]SuggestedChallenge, 1)

	for key, value := range user.Workouts {
		for _, v := range value.Results {
			var chlng SuggestedChallenge
			chlng.Type = key
			chlng.Value = v.Value
			chlng.Date = v.Date
			suggestedChlng = append(suggestedChlng, chlng)
		}
	}
	sort.Slice(suggestedChlng, func(i, j int) bool { return suggestedChlng[i].Date > suggestedChlng[j].Date })
	return suggestedChlng, nil
}

func (c *ChallengeStore) GetAllChallenges() ([]Challenge, error) {
	collection := c.DB.Database("sensors").Collection("challenges")

	options := options.Find()
	filter := bson.M{}

	var results []Challenge

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem Challenge
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, nil
}
