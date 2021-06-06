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
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Phone      string  `json:"phone"`
	OwnerName  string  `json:"owner_name"`
	Address    string  `json:"address"`
	PetName    string  `json:"pet_name"`
	PetAge     string  `json:"pet_age"`
	Token      string  `json:"token" sql:"-"`
	PetType    uint    `json:"pet"`
	BreedType  uint    `json:"breed"`
	Sex        string  `json:"sex"`
	Status     string  `json:"status"`
	Height     float64 `json:"height"`
	Weight     float64 `json:"weight"`
	AvatarFile string  `json:"avatar_file"`
}

type QrUser struct {
	gorm.Model
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	OwnerName string  `json:"owner_name"`
	Address   string  `json:"address"`
	PetName   string  `json:"pet_name"`
	PetAge    string  `json:"pet_age"`
	PetType   string  `json:"pet"`
	BreedType string  `json:"breed"`
	Sex       string  `json:"sex"`
	Status    string  `json:"status"`
	Height    float64 `json:"height"`
	Weight    float64 `json:"weight"`
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

type PetProfile struct {
	PetName   string  `json:"pet_name"`
	PetType   uint    `json:"pet_type"`
	BreedType uint    `json:"breed_type"`
	Sex       string  `json:"sex"`
	Status    string  `json:"status"`
	Height    float64 `json:"height"`
	Weight    float64 `json:"weight"`
}

type UserProfile struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type Password struct {
	SimplePassword string `json:"password"`
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

func TransformAccountToPetProfile(acc *Account) *PetProfile {

	petProfile := &PetProfile{}

	petProfile.PetName, petProfile.PetType, petProfile.BreedType, petProfile.Sex, petProfile.Status, petProfile.Height, petProfile.Weight =
		acc.PetName, acc.PetType, acc.BreedType, acc.Sex, acc.Status, acc.Height, acc.Weight

	return petProfile
}

func TransformAccountToUserProfile(acc *Account) *UserProfile {

	userProfile := &UserProfile{}

	userProfile.Name, userProfile.Email, userProfile.Phone, userProfile.Address =
		acc.OwnerName, acc.Email, acc.Phone, acc.Address

	return userProfile
}

func GetUserProfile(id uint) map[string]interface{} {

	account := &Account{}

	err := GetDB().Table("accounts").Where("id = ?", id).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	userProfile := TransformAccountToUserProfile(account)

	resp := u.Message(true, "User profile")
	resp["data"] = userProfile
	return resp
}

func GetPetProfile(id uint) map[string]interface{} {

	account := &Account{}

	err := GetDB().Table("accounts").Where("id = ?", id).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	petProfile := TransformAccountToPetProfile(account)

	resp := u.Message(true, "Pet profile")
	resp["data"] = petProfile
	return resp
}

func (userProfile *UserProfile) Edit(userId uint) map[string]interface{} {

	account := &Account{}

	err := GetDB().Table("accounts").Where("id = ?", userId).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	account.OwnerName, account.Email, account.Phone, account.Address =
		userProfile.Name, userProfile.Email, userProfile.Phone, userProfile.Address

	err = GetDB().Save(account).Error
	if err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	return u.Message(true, "Success")
}

func (petProfile *PetProfile) Edit(userId uint) map[string]interface{} {

	account := &Account{}

	err := GetDB().Table("accounts").Where("id = ?", userId).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	account.PetName, account.PetType, account.BreedType, account.Sex, account.Status, account.Height, account.Weight =
		petProfile.PetName, petProfile.PetType, petProfile.BreedType, petProfile.Sex, petProfile.Status, petProfile.Height, petProfile.Weight

	err = GetDB().Save(account).Error
	if err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	return u.Message(true, "Success")
}

func (password Password) Edit(userId uint) map[string]interface{} {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password.SimplePassword), bcrypt.DefaultCost)
	stringHashedPassword := string(hashedPassword)

	account := &Account{}
	err := GetDB().Table("accounts").Where("id = ?", userId).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	account.Password = stringHashedPassword

	err = GetDB().Save(account).Error
	if err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	return u.Message(true, "Success")
}

func GetQrUser(userId int) map[string]interface{} {

	account := &Account{}
	qrUser := &QrUser{}

	petType := &Pet{}
	petBreed := &Breed{}

	err := GetDB().Table("accounts").Where("id= ?", userId).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = GetDB().Table("pets").Where("id= ?", account.PetType).First(petType).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Pet type not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = GetDB().Table("breeds").Where("id=?", account.BreedType).First(petBreed).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	qrUser.Address, qrUser.PetType, qrUser.BreedType, qrUser.PetName, qrUser.Phone, qrUser.Email, qrUser.OwnerName,
		qrUser.Sex, qrUser.Status, qrUser.Weight, qrUser.Height, qrUser.PetAge, qrUser.ID, qrUser.CreatedAt, qrUser.UpdatedAt, qrUser.DeletedAt =
		account.Address, petType.Name, petBreed.BreedName, account.PetName, account.Phone, account.Email, account.OwnerName, account.Sex,
		account.Status, account.Weight, account.Height, account.PetAge, account.ID, account.CreatedAt, account.UpdatedAt, account.DeletedAt

	resp := u.Message(true, "Qr Profile")
	resp["data"] = qrUser
	return resp
}
