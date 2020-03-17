package handlers

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"../config"
)

type userRouter struct {
	cfg *config.Cfg
}

func NewUserRouter(cfg *config.Cfg)*userRouter{
	return &userRouter{cfg: cfg}
}

func (h *userRouter)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/users" {
		switch r.Method {
		case http.MethodGet:
			usersGetAll(w, r)
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
