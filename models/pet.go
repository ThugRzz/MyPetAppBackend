package models

import (
	u "diplomaProject/utils"
	"github.com/jinzhu/gorm"
)

type Breed struct {
	gorm.Model
	BreedName string `json:"breed_name"`
	PetType   uint   `json:"pet_type"`
}

type Pet struct {
	gorm.Model
	Name string `json:"pet_type"`
}

func GetAllPetTypes() map[string]interface{} {
	petTypes := make([]*Pet, 0)
	err := GetDB().Find(&petTypes).Error
	if err != nil {
		return u.Message(false, "Connection error. Retry")
	}

	resp := u.Message(true, "Success")
	resp["data"] = petTypes
	return resp
}

func GetAllBreeds() map[string]interface{} {
	breeds := make([]*Breed, 0)
	err := GetDB().Find(&breeds).Error
	if err != nil {
		return u.Message(false, "Connection error. Retry")
	}

	resp := u.Message(true, "Success")
	resp["data"] = breeds
	return resp
}
