package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"int20h-back-end/controllers"
	"int20h-back-end/driver"
	"int20h-back-end/models"
	"log"
	"net/http"
	"os"
)

func init() {
	gotenv.Load()
}

func main() {
	clientMongo := driver.ConnectDB(os.Getenv("MONGO_URL"))

	storeData := models.DataStore{DB: clientMongo}
	storeUser := models.UserStore{DB: clientMongo}
	storeSensor := models.SensorStore{DB: clientMongo}
	storeChallenge := models.ChallengeStore{DB: clientMongo}

	controller := controllers.Controller{
		DataStore:      &storeData,
		UserStore:      &storeUser,
		SensorStore:    &storeSensor,
		ChallengeStore: &storeChallenge,
	}
	router := mux.NewRouter()

	// router.HandleFunc("/send", controller.SendDataHandler()).Methods("POST")
	// router.HandleFunc("/send", controller.RegisterUserHandler()).Methods("POST")

	// router.HandleFunc("/get/{id}", controller.GetDataHandler()).Methods("GET")
	router.HandleFunc("/get", controller.GetAllDataHandler()).Methods("GET")
	router.HandleFunc("/register", controller.RegisterUserHandler()).Methods("POST")
	router.HandleFunc("/sensor/register", controller.RegisterSensorHandler()).Methods("POST")

	router.HandleFunc("/users", controller.GetAllUsersHandler()).Methods("GET")

	router.HandleFunc("/user/{id}", controller.GetUserHandler()).Methods("GET")
	router.HandleFunc("/exercise/save", controller.SaveOneExerciseHandler()).Methods("POST")
	router.HandleFunc("/exercise/flag/save", controller.SaveExerciseWithSensorHandler()).Methods("POST")

	router.HandleFunc("/sensor/register", controller.RegisterSensorHandler()).Methods("POST")
	router.HandleFunc("/sensor/exercise/save", controller.SaveOneSensorExerciseHandler()).Methods("POST")

	router.HandleFunc("/challenges/suggested/{id}", controller.GetAllSuggestedChallengesHandler()).Methods("GET")
	router.HandleFunc("/challanges", controller.GetAllChallengesHandler()).Methods("GET")
	router.HandleFunc("/challange/create", controller.CreateChallengeHandler()).Methods("POST")

	fmt.Println("Server is listening...")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), loggedRouter))
}
