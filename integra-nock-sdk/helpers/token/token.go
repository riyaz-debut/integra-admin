package token

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//generate jwt token from the user credentials
func GenerateToken(user_id uint, org_id uint) (string, error) {
	log.Println("inside generate token userid is :", user_id)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["org_id"] = org_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(12)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("before return the token: ", token)

	signedToken, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	log.Println("return the signedToken: ", signedToken)

	return signedToken, err

}

//check token validation
func TokenValid(receivedToken string) (string, error) {
	log.Println("Inside token valid function of admin side", receivedToken)

	token := receivedToken
	log.Println("token receive inside token valid at admin", token)

	log.Println("token in admin token valid check: ", token)
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		log.Println("error in token validationnnn checkkkk", err)
		return "", err
	}
	return token, nil
}

//extract user id from the token
func ExtractTokenID(body string) (uint, uint, error) {
	log.Println("inside extract tokenID fx of admin side", body)
	token, err := jwt.Parse(body, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		log.Println("error in tokn parsing in exractToken Id fx", err)
		return 0, 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, 0, err
		}
		orgId, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["org_id"]), 10, 32)
		if err != nil {
			return 0, 0, err
		}
		return uint(uid), uint(orgId), nil
	}
	return 0, 0, nil
}

//extract particular token
func ExtractToken(bearerToken string) string {
	log.Println("inside extract token fx of client side", bearerToken)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
