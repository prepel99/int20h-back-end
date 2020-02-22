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

	controller := controllers.Controller{
		DataStore: &storeData,
		UserStore: &storeUser,
	}
	router := mux.NewRouter()

	router.HandleFunc("/send", controller.SendDataHandler()).Methods("POST")
	router.HandleFunc("/send", controller.RegisterUserHandler()).Methods("POST")

	// router.HandleFunc("/get/{id}", controller.GetDataHandler()).Methods("GET")
	router.HandleFunc("/get", controller.GetAllDataHandler()).Methods("GET")
	router.HandleFunc("/register", controller.RegisterUserHandler()).Methods("POST")
	router.HandleFunc("/users", controller.GetAllUsersHandler()).Methods("POST")

	router.HandleFunc("/user/{id}", controller.GetUserHandler()).Methods("GET")
	router.HandleFunc("/saveexercise", controller.SaveOneExerciseHandler()).Methods("POST")

	fmt.Println("Server is listening...")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), loggedRouter))
}
