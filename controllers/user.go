package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"int20h-back-end/models"
	"io/ioutil"
	"net/http"
)

func (c *Controller) RegisterUserHandler() http.HandlerFunc {
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
		id, err := c.UserStore.RegisterUser(requestData)
		if err != nil {
			// mongoResponce.Error = err
			// json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(id)
	}
}

func (c *Controller) GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID := vars["id"]

		user, err := c.UserStore.GetUser(ID)
		if err != nil {
			// mongoResponce.Error = err
			// json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func (c *Controller) SaveOneExerciseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logr.LogErr(err)
			return
		}
		requestData := struct {
			UserID   string `json:"id"`
			Exercise models.WorkOut
		}{}

		if err := json.Unmarshal(body, &requestData); err != nil {
			logr.LogErr(err)
			return
		}
		fmt.Println(requestData)
		user, err := c.UserStore.SaveOneExercise(requestData.UserID, requestData.Exercise)
		if err != nil {
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func (c *Controller) GetAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mongoResponce, err := c.DataStore.GetAllData()
		if err != nil {
			logr.LogErr(err)
			return
		}
		if err != nil {
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(mongoResponce)
	}
}
