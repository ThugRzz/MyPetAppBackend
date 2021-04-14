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
