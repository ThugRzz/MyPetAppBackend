package models

import (
	u "diplomaProject/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

/*
Структура прав доступа JWT
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//структура для учётной записи пользователя
type Account struct {
	gorm.Model
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	OwnerName string `json:"owner_name"`
	PetName   string `json:"pet_name"`
	PetAge    string `json:"pet_age"`
	Token     string `json:"token" sql:"-"`
	PetType   uint   `json:"pet"`
	BreedType uint   `json:"breed"`
}

type User struct {
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	OwnerName     string `json:"owner_name"`
	PetName       string `json:"pet_name"`
	PetAge        string `json:"pet_age"`
	Token         string `json:"token" sql:"-"`
	PetTypeName   string `json:"pet_type" sql:"-"`
	BreedTypeName string `json:"breed_type" sql:"-"`
}

//Проверить входящие данные пользователя ...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	if account.PetName == "" {
		return u.Message(false, "Pet's name shouldn't be empty"), false
	}

	//Email должен быть уникальным
	temp := &Account{}

	//проверка на наличие ошибок и дубликатов электронных писем
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Создать новый токен JWT для новой зарегистрированной учётной записи
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //удалить пароль

	response := u.Message(true, "Account has been created")
	response["data"] = account
	return response
}

func Login(email, password string) map[string]interface{} {

	account := &Account{}
	breed := &Breed{}
	pet := &Pet{}

	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Пароль не совпадает!!
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	user := TransformAccountToUser(account)

	err = GetDB().Table("breeds").Select("breed_name").Where("ID = ?", account.BreedType).Scan(breed).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Pet breed not found")
		}
		return u.Message(false, "Connection error.  retry")
	}

	err = GetDB().Table("pets").Select("name").Where("ID = ?", account.PetType).First(pet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Pet not found")
		}
		return u.Message(false, "Connection error.  retry")
	}

	user.BreedTypeName, user.PetTypeName = breed.BreedName, pet.Name

	//Создать токен JWT
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString // Сохраните токен в ответе

	resp := u.Message(true, "Logged In")

	resp["data"] = user
	return resp
}

func TransformAccountToUser(acc *Account) *User {

	user := &User{}
	user.Email, user.Phone, user.Token, user.OwnerName = acc.Email, acc.Phone, acc.Token, acc.OwnerName
	user.PetName, user.PetAge = acc.PetName, acc.PetAge

	return user
}

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //Пользователь не найден!
		return nil
	}

	acc.Password = ""
	return acc
}
