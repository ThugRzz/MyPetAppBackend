package main

import (
	"diplomaProject/app"
	"diplomaProject/controllers"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	handleAuth(router)
	handleReferences(router)
	handlePhotoUpload(router)
	handlePet(router)
	handleProfile(router)
	handleQrUser(router)

	router.HandleFunc("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "./static/index.html")
		http.ServeFile(writer, request, "./static/script.js")
		http.ServeFile(writer, request, "./static/index.css")
	})

	fmt.Println(port)

	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("static").HTTPBox()))
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":"+port, router)

	if err != nil {
		fmt.Print(err)
	}
}

func handleAuth(router *mux.Router) {
	router.HandleFunc("/api/user/new",
		controllers.CreateAccount).Methods("POST")

	router.HandleFunc("/api/user/login",
		controllers.Authenticate).Methods("POST")
}

func handleProfile(router *mux.Router) {
	router.HandleFunc("/api/user/pet",
		controllers.PetProfile).Methods("GET")

	router.HandleFunc("/api/user/pet/edit",
		controllers.EditPetProfile).Methods("POST")

	router.HandleFunc("/api/user/profile",
		controllers.UserProfile).Methods("GET")

	router.HandleFunc("/api/user/profile/edit",
		controllers.EditUserProfile).Methods("POST")

	router.HandleFunc("/api/user/profile/password",
		controllers.EditPassword).Methods("POST")
}

func handlePhotoUpload(router *mux.Router) {
	router.HandleFunc("/api/user/avatar/upload", controllers.Handler).Methods("POST")
	router.HandleFunc("/api/user/avatar",
		controllers.GetAvatar).Methods("GET")
}

func handleReferences(router *mux.Router) {
	router.HandleFunc("/api/reference/food/{id}",
		controllers.GetFoodReferenceForBreed).Methods("GET")
	router.HandleFunc("/api/reference/foods",
		controllers.GetFoodReference).Methods("GET")

	router.HandleFunc("/api/reference/care/{id}",
		controllers.GetCareReferenceForBreed).Methods("GET")
	router.HandleFunc("/api/reference/cares",
		controllers.GetCareReference).Methods("GET")

	router.HandleFunc("/api/reference/disease/{id}",
		controllers.GetDiseaseReferenceForBreed).Methods("GET")
	router.HandleFunc("/api/reference/diseases",
		controllers.GetDiseaseReference).Methods("GET")

	router.HandleFunc("/api/reference/training/{id}",
		controllers.GetTrainingReferenceForBreed).Methods("GET")
	router.HandleFunc("/api/reference/trainings",
		controllers.GetTrainingReference).Methods("GET")
}

func handlePet(router *mux.Router) {
	router.HandleFunc("/api/pet/types",
		controllers.GetPetTypes).Methods("GET")

	router.HandleFunc("/api/pet/breeds",
		controllers.GetBreeds).Methods("GET")
}

func handleQrUser(router *mux.Router) {
	router.HandleFunc("/api/user/qr/{id}",
		controllers.QrProfile).Methods("GET")
}
