package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SensorValue struct {
	Value int   `json:"value" bson:"value"`
	Date  int64 `json:"date" bson:"date"`
}

type SensorData struct {
	ID         string        `json:"id" bson:"_id"`
	Values     []SensorValue `json:"values" bson:"values"`
	Updated_at int64         `json:"updated_at" bson:"updated_at"`
}

type RequestData struct {
	Data []SensorData `json:"data"`
}

type ResponceSendData struct {
	CreatedIDs []string `json:"createdIDs"`
	UpdatedIDs []string `json:"updatedIDs"`
	Error      error    `json:"error"`
}

type DataStorer interface {
	SaveData([]SensorData) (ResponceSendData, error)
	GetAllData() ([]SensorData, error)
}

type DataStore struct {
	DB *mongo.Client
}

func (d *DataStore) SaveData(data []SensorData) (ResponceSendData, error) {
	collection := d.DB.Database("sensors").Collection("SensorsData")

	var updatedIDs = make([]interface{}, 0) // existed elements were updated
	var createdIDs = make([]interface{}, 0) // new elements were inserted

	var sensor SensorData
	var sensorFound bool

	for _, element := range data {
		now := time.Now().Unix()
		element.Updated_at = now

		filter := bson.D{{"_id", element.ID}}
		sensor = SensorData{}
		sensorFound = true

		err := collection.FindOne(context.TODO(), filter).Decode(&sensor)
		if err != nil {
			if err.Error() != "mongo: no documents in result" {
				return ResponceSendData{Error: err}, err
			}
			if err.Error() == "mongo: no documents in result" {
				sensor = element
				sensorFound = false
			}
		}

		if sensorFound {
			for _, v := range element.Values {
				if v.Date == 0 {
					v.Date = time.Now().Unix()
				}
				sensor.Values = append(sensor.Values, v)
			}
		} else {
			for _, v := range sensor.Values {
				if v.Date == 0 {
					v.Date = time.Now().Unix()
				}
			}
		}

		doc, err := toDoc(sensor)
		if err != nil {
			return ResponceSendData{Error: err}, err
		}
		update := bson.D{{"$set", *doc}}
		result, err := collection.UpdateOne(
			context.Background(),
			filter,
			update,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return ResponceSendData{Error: err}, err
		}

		if result.ModifiedCount == 1 {
			updatedIDs = append(updatedIDs, element.ID)
		}
		if result.UpsertedID != nil {
			createdIDs = append(createdIDs, result.UpsertedID)
		}
	}

	// convert []interface to []string
	resultUpdatedIDs := make([]string, len(updatedIDs))
	for i, v := range updatedIDs {
		resultUpdatedIDs[i] = fmt.Sprint(v)
	}

	resultCreatedIDs := make([]string, len(createdIDs))
	for i, v := range createdIDs {
		resultCreatedIDs[i] = fmt.Sprint(v)
	}

	result := ResponceSendData{
		CreatedIDs: resultCreatedIDs,
		UpdatedIDs: resultUpdatedIDs,
	}

	return result, nil
}

func (d *DataStore) GetAllData() ([]SensorData, error) {
	collection := d.DB.Database("sensors").Collection("SensorsData")

	options := options.Find()
	filter := bson.M{}

	// Here's an array in which you can store the decoded documents
	var results []SensorData

	// Passing nil as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem SensorData
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, nil
}

// func (d *DataStore) GetData(string) ([]SensorData, error) {
// 	collection := d.DB.Database("sensors").Collection("SensorsData")

// 	options := options.Find()
// 	filter := bson.M{}

// 	// Here's an array in which you can store the decoded documents
// 	var results []SensorData

// 	// Passing nil as the filter matches all documents in the collection
// 	cur, err := collection.Find(context.TODO(), filter, options)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Finding multiple documents returns a cursor
// 	// Iterating through the cursor allows us to decode documents one at a time
// 	for cur.Next(context.TODO()) {
// 		// create a value into which the single document can be decoded
// 		var elem SensorData
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, elem)
// 	}
// 	return results, nil
// }

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}
