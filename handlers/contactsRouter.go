package handlers

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func ContactRouter(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		postError(w, http.StatusNotFound)
		return
	}
	objectID, _ := primitive.ObjectIDFromHex(id)

	switch r.Method {
	case http.MethodGet:
		contactsGet(w, r, objectID)
		return
	case http.MethodPost:
		contactPostOne(w, r)
		return
	case http.MethodPut:
		contactPutOne(w, r, objectID)
		return
	case http.MethodDelete:
		contactDeleteOne(w, r, objectID)
		return
	case http.MethodOptions:
		contactsDeleteByOwner(w, r, id)
		return
	default:
		postError(w, http.StatusMethodNotAllowed)
		return
	}

}
