package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/1994benc/generic-auth-service/token"
	"github.com/go-oauth2/oauth2/utils/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestHandleVerifyInvalidToken(t *testing.T) {
	mockJWTSecret := "asjcascjacsplpcpa[sop[coapscoosaicias0c0-as-coa0iscoasjckjasok"
	maker, err := token.NewJWTMaker(mockJWTSecret)
	require.NoError(t, err)
	require.NotNil(t, maker)

	mockDb := &mongo.Client{}

	req := httptest.NewRequest(http.MethodGet, "/verify-token?application_access_token=blahblah", nil)
	w := httptest.NewRecorder()
	HandleVerifyToken(mockDb, maker)(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(data), "access token is invalid")
	require.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestHandleVerifyValidToken(t *testing.T) {
	mockJWTSecret := "asjcascjacsplpcpa[sop[coapscoosaicias0c0-as-coa0iscoasjckjasok"
	maker, err := token.NewJWTMaker(mockJWTSecret)
	require.NoError(t, err)
	require.NotNil(t, maker)

	userId, err := uuid.NewRandom()
	require.NoError(t, err)

	token, err := maker.CreateToken(userId.String(), time.Minute)
	require.NoError(t, err)
	require.NotNil(t, token)

	mockDb := &mongo.Client{}

	req := httptest.NewRequest(http.MethodGet, "/verify-token?application_access_token="+token, nil)
	w := httptest.NewRecorder()
	HandleVerifyToken(mockDb, maker)(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(data), "user_id")
	require.Equal(t, http.StatusOK, res.StatusCode)
}
