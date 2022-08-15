package main

import (
	"context"
	"net/http"
	"os"

	"github.com/1994benc/generic-auth-service/config"
	"github.com/1994benc/generic-auth-service/handlers"
	"github.com/1994benc/generic-auth-service/token"
	"github.com/joho/godotenv"
)

func main() {
	// if .env.local file exists, load it
	if _, err := os.Stat(".env.local"); err == nil {
		godotenv.Load(".env.local")
	}

	db := config.ConnectToDb()
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	secret_key := os.Getenv("JWT_SECRET")
	if secret_key == "" {
		panic("JWT_SECRET is not set")
	}
	token_maker, err := token.NewJWTMaker(secret_key)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handlers.HandleHome(db))

	// provides this endpoint with an oauth access token
	// if the oauth access token is valid,
	// it will return a token that can be used to access our resources
	http.HandleFunc("/exchange-oauth-token-for-application-token", handlers.HandleExchangeToken(db, token_maker))

	// verifies that the application access token is valid
	// this is usually used by resource servers to verify that the application access token is valid
	http.HandleFunc("/verify-token", handlers.HandleVerifyToken(db, token_maker))

	// refreshes the application access token.
	// call this endpoint with the refresh token to get a new application access token
	// and also a new refresh token
	http.HandleFunc("/refresh-token", handlers.HandleRefreshToken(db, token_maker))

	http.HandleFunc("/revoke-refresh-token", handlers.HandleRevokeRefreshToken)

	http.ListenAndServe(":8080", nil)
}
