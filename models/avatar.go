package models

import (
	u "diplomaProject/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	"time"
)

type Avatar struct {
	Url string `json:"avatar_url"`
}

func SaveAvatar(userId uint, avatarUrl string, tempFileName string) map[string]interface{} {

	account := &Account{}

	avatar := &Avatar{}

	err := GetDB().Table("accounts").Where("id = ?", userId).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	account.AvatarFile = tempFileName
	avatar.Url = avatarUrl

	err = GetDB().Save(account).Error
	if err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	resp := u.Message(true, "Success")
	resp["data"] = avatar

	return resp
}

func GetAvatar(userId uint, s *session.Session) map[string]interface{} {

	account := &Account{}

	avatar := Avatar{}

	err := GetDB().Table("accounts").Where("id = ?", userId).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	fileName := account.AvatarFile

	if fileName == "" {
		return u.Message(false, "Avatar not found")
	}

	svc := s3.New(s)

	param := &s3.GetObjectInput{
		Bucket: aws.String("mypetapp"),
		Key:    aws.String(fileName),
	}

	req, _ := svc.GetObjectRequest(param)

	url, err := req.Presign(60 * time.Minute)
	if err != nil {
		return u.Message(false, "Error")
	}

	avatar.Url = url

	resp := u.Message(true, "Avatar")
	resp["data"] = avatar
	return resp
}
