package controllers

import (
	"bytes"
	"diplomaProject/models"
	u "diplomaProject/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader,
	w http.ResponseWriter, r *http.Request) {
	// get the file size and read
	// the file content into a buffer

	tokenHeader := r.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]

	tk := &models.Token{}
	_, _ = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := "pictures/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	svc := s3.New(s)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("mypetapp"),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		u.Respond(w, u.Message(false, "Error"))
		return
	}
	param := &s3.GetObjectInput{
		Bucket: aws.String("mypetapp"),
		Key:    aws.String(tempFileName),
	}

	req, _ := svc.GetObjectRequest(param)

	url, err := req.Presign(60 * time.Minute)

	if err != nil {
		u.Respond(w, u.Message(false, "Error"))
		return
	}

	resp := models.SaveAvatar(tk.UserId, url, tempFileName)
	u.Respond(w, resp)
}

func GetAvatar(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]

	tk := &models.Token{}
	_, _ = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	region := os.Getenv("region")
	secretId := os.Getenv("s3_id")
	secretKey := os.Getenv("s3_secret")
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(secretId, secretKey, ""),
	})
	if err != nil {
		u.Respond(w, u.Message(false, "Could not upload file aws"))
		return
	}

	resp := models.GetAvatar(tk.UserId, s)
	u.Respond(w, resp)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	maxSize := int64(102400000) // allow only 1MB of file size

	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		u.Respond(w, u.Message(false, "Image too large"))
		return
	}

	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		u.Respond(w, u.Message(false, "Could not get uploaded file"))
		return
	}
	defer file.Close()

	region := os.Getenv("region")
	secretId := os.Getenv("s3_id")
	secretKey := os.Getenv("s3_secret")
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(secretId, secretKey, ""),
	})
	if err != nil {
		u.Respond(w, u.Message(false, "Could not upload file aws"))
		return
	}

	UploadFileToS3(s, file, fileHeader, w, r)
}
