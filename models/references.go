package models

import "github.com/jinzhu/gorm"

type Food struct {
	gorm.Model
	BreedType   *Breed `json:"breed_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Care struct {
	gorm.Model
	BreedType   *Breed `json:"breed_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Disease struct {
	gorm.Model
	BreedType   *Breed `json:"breed_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Training struct {
	gorm.Model
	BreedType   *Breed `json:"breed_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
