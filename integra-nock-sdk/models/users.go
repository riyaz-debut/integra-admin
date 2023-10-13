package models

import (
	"errors"
	"integra-nock-sdk/database"
	"integra-nock-sdk/helpers/token"
	"log"
	"reflect"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        uint      `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"user_name" binding:"required" `
	Password  string    `json:"password" binding:"required"`
	OrgId     int       `json:"org_id,omitempty"`
	OrgMsp    string    `json:"org_msp,omitempty"`
	Status	  int		`json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

//check username and password is match or not
func AdminCredentialsCheck(username string, password string) (string, error) {

	var err error
	var userName = username

	log.Println("username is:", userName)
	user := User{}

	err = database.Connector.Model(User{}).Where("user_name = ? AND role = ?", username, "admin").Take(&user).Error
	if err != nil {
		return "", err
	}

	//verify password is match or not
	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("password check error", err)
		return "", nil
	}

	log.Println("userid", user.Id)
	log.Println("org id ", user.OrgId)
	log.Println(reflect.TypeOf(user.Id))
	token, err := token.GenerateToken(user.Id, uint(user.OrgId))
	if err != nil {
		log.Println("error in genrate token function", err)
		return "", nil
	}
	return token, nil

}

func UserCredentialsCheck(username string, password string) (string, error) {

	var err error
	var userName = username

	log.Println("username is:", userName)
	user := User{}
	/////////////////////

	// if err := database.Connector.Table("users").Find(&user, "user_name = ?", userName).Error; err != nil {

	// 	return "", err

	// }

	// user_id := user.Id
	// log.Println("user id in usercheckkk controller: ", user_id)
	// log.Println("user data in usercheckkk controller: ", user)

	// status := user.Status
	// log.Println("user status in disable controller: ", status)

	// // check for status for disabling the user
	// if  status == 0 {
	// 	return "", err
	// } 

	////////////////////
	// err = database.Connector.Model(User{}).Where("user_name = ? AND role = ?", username, "user").Take(&user).Error
	err = database.Connector.Model(User{}).Where("user_name = ?", username).Take(&user).Error
	if err != nil {
		return "", err
	}

	//verify password is match or not
	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("password check error", err)
		return "", nil
	}

	log.Println("userid", user.Id)
	log.Println("org id ", user.OrgId)
	log.Println(reflect.TypeOf(user.Id))
	token, err := token.GenerateToken(user.Id, uint(user.OrgId))
	if err != nil {
		log.Println("error in genrate token function", err)
		return "", nil
	}
	return token, nil

}

//get particular user by user id
func GetUserByID(uid uint) (User, error) {
	var user User
	if err := database.Connector.First(&user, uid).Error; err != nil {
		return user, errors.New("User with this id not found!")
	}
	log.Println("user data found ", user)
	//user.PrepareGive()

	return user, nil
}
