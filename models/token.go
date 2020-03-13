package models

import "github.com/dgrijalva/jwt-go"

type Token struct {
	Name string `json:"name"`
	jwt.StandardClaims
}
