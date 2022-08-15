package handlers

import "net/http"

func HandleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("revoke-refresh-token"))
}
