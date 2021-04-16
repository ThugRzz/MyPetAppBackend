package main

import (
	"diplomaProject/app"
	"diplomaProject/controllers"
	"fmt"
	"github.com/gorilla/mux"
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

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)

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
