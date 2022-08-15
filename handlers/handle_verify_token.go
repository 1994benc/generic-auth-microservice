package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/1994benc/generic-auth-service/common"
	"github.com/1994benc/generic-auth-service/token"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleVerifyToken(db *mongo.Client, token_maker token.Maker) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		// get application access token from query string
		applicationAccessToken := r.FormValue("application_access_token")
		if applicationAccessToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "application_access_token is required"})
			if err != nil {
				panic(err)
			}
		}

		payload, err := token_maker.ValidateToken(applicationAccessToken)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "access token is invalid"})
			if err != nil {
				panic(err)
			}
			return
		}

		// return payload if everything is ok
		err = json.NewEncoder(w).Encode(payload)
		if err != nil {
			err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "failed to write and send payload"})
			if err != nil {
				panic(err)
			}
			return
		}
	}
	return fn
}
