package models

import (
	u "diplomaProject/utils"
	"github.com/jinzhu/gorm"
)

type Food struct {
	gorm.Model
	BreedType   uint   `json:"breed_type"`
	PetType     uint   `json:"pet_type"`
	Title       string `json:"title"`
	StepTitle   string `json:"step_title"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

type Care struct {
	gorm.Model
	BreedType   uint   `json:"breed_type"`
	PetType     uint   `json:"pet_type"`
	Title       string `json:"title"`
	StepTitle   string `json:"step_title"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

type Disease struct {
	gorm.Model
	BreedType   uint   `json:"breed_type"`
	PetType     uint   `json:"pet_type"`
	Title       string `json:"title"`
	StepTitle   string `json:"step_title"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

type Training struct {
	gorm.Model
	BreedType   uint   `json:"breed_type"`
	PetType     uint   `json:"pet_type"`
	Title       string `json:"title"`
	StepTitle   string `json:"step_title"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

func GetAllFoodReference() map[string]interface{} {

	foods := make([]*Food, 0)
	err := GetDB().Find(&foods).Error
	if err != nil {
		return u.Message(false, "Connection error. Retry")
	}

	resp := u.Message(true, "Success")
	resp["data"] = foods
	return resp
}

func GetFoodReference(breedType uint) map[string]interface{} {

	food := &Food{}

	err := GetDB().Table("food").Where("breed_type = ?", breedType).First(food).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Food reference for this breed not found")
		}
		return u.Message(false, "Connection error.  retry")
	}

	resp := u.Message(true, "Success")

	resp["data"] = food
	return resp
}

func GetAllCareReference() map[string]interface{} {

	cares := make([]*Care, 0)
	err := GetDB().Find(&cares).Error
	if err != nil {
		return u.Message(false, "Connection error. Retry")
	}

	resp := u.Message(true, "Success")
	resp["data"] = cares
	return resp
}

func GetCareReference(breedType uint) map[string]interface{} {

	care := &Care{}

	err := GetDB().Table("cares").Where("breed_type = ?", breedType).First(care).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Food reference for this breed not found")
		}
		return u.Message(false, "Connection error.  retry")
	}

	resp := u.Message(true, "Success")

	resp["data"] = care
	return resp
}

func GetAllDiseaseReference() map[string]interface{} {

	diseases := make([]*Disease, 0)
	err := GetDB().Find(&diseases).Error
	if err != nil {
		return u.Message(false, "Connection error. Retry")
	}

	resp := u.Message(true, "Success")
	resp["data"] = diseases
	return resp
}

func GetDiseaseReference(breedType uint) map[string]interface{} {

	disease := &Disease{}

	err := GetDB().Table("diseases").Where("breed_type = ?", breedType).First(disease).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Food reference for this breed not found")
		}
		return u.Message(false, "Connection error.  retry")
	}

	resp := u.Message(true, "Success")

	resp["data"] = disease
	return resp
}

func GetAllTrainingReference() map[string]interface{} {

	trainings := make([]*Training, 0)
	err := GetDB().Find(&trainings).Error
	if err != nil {
		return u.Message(false, "Connection error. Retry")
	}

	resp := u.Message(true, "Success")
	resp["data"] = trainings
	return resp
}

func GetTrainingReference(breedType uint) map[string]interface{} {

	training := &Training{}

	err := GetDB().Table("trainings").Where("breed_type = ?", breedType).First(training).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Food reference for this breed not found")
		}
		return u.Message(false, "Connection error.  retry")
	}

	resp := u.Message(true, "Success")

	resp["data"] = training
	return resp
}
