package controllers

import (
	"encoding/json"
	"fmt"
	"int20h-back-end/models"
	"io/ioutil"
	"net/http"
)

func (c *Controller) RegisterSensorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logr.LogErr(err)
			return
		}

		requestData := models.User{}

		if err := json.Unmarshal(body, &requestData); err != nil {
			// .Error = err
			// json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}
		fmt.Println(requestData)
		id, err := c.SensorStore.RegisterUser(requestData)
		if err != nil {
			// mongoResponce.Error = err
			// json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(id)
	}
}

func (c *Controller) SaveOneSensorExerciseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logr.LogErr(err)
			return
		}
		requestData := struct {
			Exercise models.WorkOut
		}{}

		if err := json.Unmarshal(body, &requestData); err != nil {
			logr.LogErr(err)
			return
		}
		_, err = c.SensorStore.SaveOneExercise("5e51966d09eaf8c6d663ff3c", requestData.Exercise)
		if err != nil {
			logr.LogErr(err)
			return
		}
	}
}
