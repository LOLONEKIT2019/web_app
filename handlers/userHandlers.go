package handlers

import (
	. "../models"
	. "../storage"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func bodyToUser(r *http.Request, u *User) error {
	if r.Body == nil {
		return errors.New("requset body is impte !")
	}
	if u == nil {
		return errors.New("a storage is required")
	}
	return json.NewDecoder(r.Body).Decode(u)

}

func usersGetAll(w http.ResponseWriter, _ *http.Request) {

	users, err := FindAllUsers()
	if err != nil {
		log.Println(err)
		return
	}

	postBodyResponse(w, http.StatusOK, jsonResponse{
		"users": users,
	})
	return
}

//create new user
func usersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	if !u.CheckMail() {
		postBodyResponse(w, http.StatusBadRequest, jsonResponse{"mail not correct": u.Mail})
		return
	}

	if !u.CheckPassword() {
		postBodyResponse(w, http.StatusBadRequest, jsonResponse{"Problems w password length (need 8 - 16 symbils not": len(u.Password)})
		return
	}

	err = CrateUser(u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", "/users/"+u.Mail)
	w.WriteHeader(http.StatusCreated)

}

func usersGetOne(w http.ResponseWriter, _ *http.Request, id primitive.ObjectID) {
	u, err := FindOneUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"storage": u})
}

//update user pass || status
func usersPutOne(w http.ResponseWriter, r *http.Request, id primitive.ObjectID) {
	u := new(User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}

	if len(u.Status) != 0 {
		if u.CheckStatus() {
			postError(w, http.StatusBadRequest)
			return
		}
	}

	if len(u.Password) != 0 {
		if u.CheckPassword() {
			postError(w, http.StatusBadRequest)
			return
		}
	}

	u.Modified = time.Now()

	err = UpdateUser(id, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func usersDeleteOne(w http.ResponseWriter, _ *http.Request, id primitive.ObjectID) {

	err := DeleteOneUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
