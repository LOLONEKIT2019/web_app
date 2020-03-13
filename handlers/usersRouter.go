package handlers

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

func UsersRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/users" {
		switch r.Method {
		case http.MethodGet:
			UsersGetAll(w, r)
			return
		case http.MethodPost:
			usersPostOne(w, r)
			return
		default:
			postError(w, http.StatusMethodNotAllowed)
			return
		}
	}

	path = strings.TrimPrefix(path, "/users/")
	if path == "" {
		postError(w, http.StatusNotFound)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	objectID, _ := primitive.ObjectIDFromHex(id)

	switch r.Method {
	case http.MethodGet:
		usersGetOne(w, r, objectID)
		return
	case http.MethodPut:
		usersPutOne(w, r, objectID)
		return
	case http.MethodDelete:
		usersDeleteOne(w, r, objectID)
		return
	default:
		postError(w, http.StatusMethodNotAllowed)
		return
	}

}
