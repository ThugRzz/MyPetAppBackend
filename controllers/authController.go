package controllers

import (
	"diplomaProject/models"
	u "diplomaProject/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create() //Создать аккаунт
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

var UserProfile = func(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]

	tk := &models.Token{}
	_, _ = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	resp := models.GetUserProfile(tk.UserId)
	u.Respond(w, resp)
}

var PetProfile = func(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]

	tk := &models.Token{}
	_, _ = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	resp := models.GetPetProfile(tk.UserId)
	u.Respond(w, resp)
}

var EditUserProfile = func(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]

	tk := &models.Token{}
	_, _ = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	userProfile := &models.UserProfile{}
	err := json.NewDecoder(r.Body).Decode(userProfile)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := userProfile.Edit(tk.UserId)
	u.Respond(w, resp)
}

var EditPetProfile = func(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]

	tk := &models.Token{}
	_, _ = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	petProfile := &models.PetProfile{}
	err := json.NewDecoder(r.Body).Decode(petProfile) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := petProfile.Edit(tk.UserId) //Создать аккаунт
	u.Respond(w, resp)
}

var EditPassword = func(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]

	tk := &models.Token{}
	_, _ = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	password := &models.Password{}
	err := json.NewDecoder(r.Body).Decode(password) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := password.Edit(tk.UserId)
	u.Respond(w, resp)
}

var QrProfile = func(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	resp := models.GetQrUser(id)
	u.Respond(w, resp)
}
