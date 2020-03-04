package main

import (
	"log"
	"reflect"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = "thisissecret"

// CustomeClaims ...
// type CustomeClaims struct {
// 	Expired bool   `json:"expired"`
// 	UserID  string `json:"user_id"`
// }

// // Valid for jwt.Claims interface ...
// func (c CustomeClaims) Valid() error {
// 	return nil
// }

func main() {
	// generate
	claims := jwt.MapClaims{"expired": true, "user_id": "2312312"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("tokenString is: ", tokenString)
	// parse

	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { return []byte(secretKey), nil })
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(parsedToken)
	if !parsedToken.Valid {
		log.Println("invalid token string")
		return
	}
	parsedClaims := parsedToken.Claims.(jwt.MapClaims)
	if !reflect.DeepEqual(claims, parsedClaims) {
		log.Printf("not equal claims, want: %v, got: %v\n", claims, parsedClaims)
		return
	}
	log.Println("success")
}
