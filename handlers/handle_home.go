package handlers

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func HandleHome(db *mongo.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the Auth service.")
	}

	return http.HandlerFunc(fn)
}
