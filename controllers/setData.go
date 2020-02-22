package controllers

import (
	"encoding/json"
	"int20h-back-end/logger"
	"int20h-back-end/models"
	"io/ioutil"
	"net/http"
)

type Controller struct {
	DataStore models.DataStorer
	UserStore models.UserStorer
}

var logr logger.Logger

func (c *Controller) SendDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logr.LogErr(err)
			return
		}

		requestData := models.RequestData{}
		mongoResponce := models.ResponceSendData{}

		w.Header().Set("Content-Type", "application/json")
		if err := json.Unmarshal(body, &requestData); err != nil {
			mongoResponce.Error = err
			json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}

		mongoResponce, err = c.DataStore.SaveData(requestData.Data)
		if err != nil {
			mongoResponce.Error = err
			json.NewEncoder(w).Encode(mongoResponce)
			logr.LogErr(err)
			return
		}
		json.NewEncoder(w).Encode(mongoResponce)
	}
}

func (c *Controller) GetAllDataHandler() http.HandlerFunc {
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
