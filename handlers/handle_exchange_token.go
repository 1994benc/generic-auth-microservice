package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/1994benc/generic-auth-service/common"
	"github.com/1994benc/generic-auth-service/token"
	"github.com/1994benc/generic-auth-service/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoogleUser struct {
	Email   string `json:"email"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// exchange oauth access token for application access token.
// this handles verifying with third party provider if the access token provided is valid
// if so it will return 'our' application access token and refresh token that can be used to access our resources
func HandleExchangeToken(db *mongo.Client, token_maker token.Maker) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// get provider from query string
		provider := r.FormValue("provider")
		if provider == "" {
			// return bad request if provider is not specified
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "provider is not specified")
			return
		}

		if provider == "google" {

			var internal_user_id string
			// get access token from query string
			accessToken := r.FormValue("access_token")
			if accessToken == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "access_token is required")
				return
			}

			// call google api to verify access token
			// get google_user info from google api
			// parse response to get google_user info
			google_user, shouldReturn := getGoogleUserInfo(accessToken, w)

			if shouldReturn {
				return
			}

			if google_user.ID == "" {
				w.WriteHeader(http.StatusForbidden)
				err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "access token is invalid"})
				if err != nil {
					panic(err)
				}
				return
			}

			email := google_user.Email

			// get user from db
			user_data, err := users.GetUserByEmail(db, email)

			if user_data != nil {
				internal_user_id = user_data.ID
			}

			if user_data == nil || err != nil {
				// create user in db
				user := &users.User{
					ID:      google_user.ID,
					Email:   google_user.Email,
					Name:    google_user.Name,
					Picture: google_user.Picture,
				}

				inserted_id, err := users.CreateUser(db, user)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "Failed to create user in db"})
					if err != nil {
						panic(err)
					}
					return
				}

				internal_user_id = inserted_id
			}

			println("internal user id: " + internal_user_id)

			// generate acess token for the user that expires in 15 minutes
			access_token, err := token_maker.CreateToken(internal_user_id, time.Minute*15)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "failed to create application access token"})
				if err != nil {
					panic(err)
				}
				return
			}

			// generate refresh token for the user that expires in 1 month
			refresh_token, err := token_maker.CreateToken(internal_user_id, time.Minute*60*24*30)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "failed to create refresh token"})
				if err != nil {
					panic(err)
				}
				return
			}

			// return token as json
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(common.ApplicationAccessToken{AccessToken: access_token, RefreshToken: refresh_token})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err := json.NewEncoder(w).Encode(common.ErrorWithInfo{Message: "failed writing tokens"})
				if err != nil {
					panic(err)
				}
				return
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "provider not supported")
			return
		}
	}

	return http.HandlerFunc(fn)
}

func getGoogleUserInfo(accessToken string, w http.ResponseWriter) (GoogleUser, bool) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(accessToken))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to get user info from google api")
		return GoogleUser{}, true
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to read response body")
		return GoogleUser{}, true
	}

	var user GoogleUser
	if err := json.Unmarshal(response, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to parse response body")
		return GoogleUser{}, true
	}
	return user, false
}
