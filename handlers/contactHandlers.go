package handlers

import (
	. "../models"
	. "../storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func bodyToContact(r *http.Request, c *Contact) error {
	if r.Body == nil {
		return errors.New("requset body is impte !")
	}
	if c == nil {
		return errors.New("a storage is required")
	}
	return json.NewDecoder(r.Body).Decode(c)

}

//create new user
func contactPostOne(w http.ResponseWriter, r *http.Request) {
	c := new(Contact)
	err := bodyToContact(r, c)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	if !c.CheckPhone() {
		postBodyResponse(w, http.StatusBadRequest, jsonResponse{"Phone number not correct": c.Phone})
		return
	}

	if !c.CheckName() {
		postBodyResponse(w, http.StatusBadRequest, jsonResponse{"Problems w name length (need 3 - 16 symbils) not": len(c.Name)})
		return
	}

	vars := mux.Vars(r)
	c.Owner = vars["id"]

	c.Created = time.Now()
	c.Modified = time.Now()

	err = CrateContact(c)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func contactsGet(w http.ResponseWriter, _ *http.Request, id primitive.ObjectID) {
	c, err := FindContacts(id.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"contacts": c})
}

//update contact
func contactPutOne(w http.ResponseWriter, r *http.Request, id primitive.ObjectID) {
	c := new(Contact)
	err := bodyToContact(r, c)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}

	err = UpdateContact(id, c)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func contactDeleteOne(w http.ResponseWriter, _ *http.Request, id primitive.ObjectID) {

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

func contactsDeleteByOwner(w http.ResponseWriter, _ *http.Request, owner string) {

	err := DeleteContacts(owner)
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
