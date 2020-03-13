package handlers

import (
	"../models"
	. "../security"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const tokenKey = "MyApi"

func bodyToAuthData(r *http.Request, auth *AuthData) error {
	if r.Body == nil {
		return errors.New("requset body is impte !")
	}

	return json.NewDecoder(r.Body).Decode(auth)
}

func CreateToken(w http.ResponseWriter, r *http.Request) {

	user := new(AuthData)
	err := bodyToAuthData(r, user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ok, err := CheckUser(user)
	if err != nil || !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Token{
		Name:           user.Mail,
		StandardClaims: jwt.StandardClaims{},
	})
	jwtToken, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		postBodyResponse(w, http.StatusInternalServerError, jsonResponse{"error": "token_generation_failed"})
		return
	}

	postBodyResponse(w, http.StatusOK, jsonResponse{"jwt": jwtToken})
	return
}

func JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/token/new", "/token/login"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {

				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			postError(w, http.StatusNotFound)
			return
		}

		//basic / bearer
		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			postError(w, http.StatusForbidden)
			return
		}

		if splitted[0] != "Basic" && splitted[0] != "Bearer" {
			postError(w, http.StatusForbidden)
			return
		}

		tokenString := splitted[1] //Grab the token part, what we are truly interested in

		tk := models.Token{}

		token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenKey), nil
		})

		if err != nil {
			postError(w, http.StatusForbidden)
			return
		}

		if !token.Valid {
			postError(w, http.StatusUnauthorized)
			return
		}

		fmt.Println(fmt.Sprintf("User %v made some changes !", tk.Name)) //Useful for monitoring

		ctx := context.WithValue(r.Context(), "user", tk.Name)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
