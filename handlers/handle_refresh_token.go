package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/1994benc/generic-auth-service/common"
	"github.com/1994benc/generic-auth-service/token"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleRefreshToken(db *mongo.Client, token_maker token.Maker) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// get refresh token from query string
		refreshToken := r.FormValue("refresh_token")
		if refreshToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "refresh_token is required in query string"})
			if err != nil {
				panic(err)
			}
			return
		}

		// verify if the refresh token is valid
		payload, err := token_maker.ValidateToken(refreshToken)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "refresh token is invalid"})
			if err != nil {
				panic(err)
			}
			return
		}

		// generate new access token and refresh token
		newRefreshToken, err := token_maker.CreateToken(payload.UserID, time.Minute*60*24*30)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "failed to create new refresh token"})
			if err != nil {
				panic(err)
			}
			return
		}
		newAccessToken, err := token_maker.CreateToken(payload.UserID, time.Minute*15)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "failed to create new access token"})
			if err != nil {
				panic(err)
			}
			return
		}

		// return new access token and refresh token
		err = json.NewEncoder(w).Encode(common.ApplicationAccessToken{AccessToken: newAccessToken, RefreshToken: newRefreshToken})
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
