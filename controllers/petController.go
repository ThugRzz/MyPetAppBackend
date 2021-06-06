package controllers

import (
	"diplomaProject/models"
	u "diplomaProject/utils"
	"net/http"
)

var GetPetTypes = func(w http.ResponseWriter, r *http.Request) {
	resp := models.GetAllPetTypes()
	u.Respond(w, resp)
}

var GetBreeds = func(w http.ResponseWriter, r *http.Request) {
	resp := models.GetAllBreeds()
	u.Respond(w, resp)
}
