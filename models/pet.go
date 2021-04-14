package models

import "github.com/jinzhu/gorm"

type Breed struct {
	gorm.Model
	BreedName string `json:"breed_name"`
	PetType   uint   `json:"pet_type"`
}

type Pet struct {
	gorm.Model
	Name string `json:"pet_type"`
}
