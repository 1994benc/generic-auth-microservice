package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetGoogleOauthConfig() *oauth2.Config {
	var (
		googleOauthConfig = &oauth2.Config{
			RedirectURL:  "http://localhost:8080/google-callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		}
	)

	return googleOauthConfig
}
