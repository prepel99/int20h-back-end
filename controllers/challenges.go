package controllers

import (
	"encoding/json"
	// "fmt"
	// "github.com/gorilla/mux"
	"int20h-back-end/models"
	"io/ioutil"
	"net/http"
)

func (c *Controller) GetAllSuggestedChallengesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		challenges, err := c.ChallengeStore.GetAllChallenges()
		if err != nil {
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(challenges)
	}
}

func (c *Controller) GetAllChallengesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mongoResponce, err := c.ChallengeStore.GetAllChallenges()
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

func (c *Controller) CreateChallengeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logr.LogErr(err)
			return
		}

		requestData := models.Challenge{}

		if err := json.Unmarshal(body, &requestData); err != nil {
			// .Error = err
			// json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}
		id, err := c.ChallengeStore.CreateChallenge(requestData)
		if err != nil {
			// mongoResponce.Error = err
			// json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(id)
	}
}
