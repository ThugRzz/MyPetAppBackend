package controllers

import (
	"diplomaProject/models"
	u "diplomaProject/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetFoodReference = func(w http.ResponseWriter, r *http.Request) {
	resp := models.GetAllFoodReference()
	u.Respond(w, resp)
}

var GetDiseaseReference = func(w http.ResponseWriter, r *http.Request) {
	resp := models.GetAllDiseaseReference()
	u.Respond(w, resp)
}

var GetTrainingReference = func(w http.ResponseWriter, r *http.Request) {
	resp := models.GetAllTrainingReference()
	u.Respond(w, resp)
}

var GetCareReference = func(w http.ResponseWriter, r *http.Request) {
	resp := models.GetAllCareReference()
	u.Respond(w, resp)
}

var GetFoodReferenceForBreed = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}
	resp := models.GetFoodReference(uint(id))

	u.Respond(w, resp)
}

var GetCareReferenceForBreed = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}
	resp := models.GetCareReference(uint(id))

	u.Respond(w, resp)
}

var GetDiseaseReferenceForBreed = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}
	resp := models.GetDiseaseReference(uint(id))

	u.Respond(w, resp)
}

var GetTrainingReferenceForBreed = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}
	resp := models.GetTrainingReference(uint(id))

	u.Respond(w, resp)
}
