package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var api *API

func TestMain(m *testing.M) {
	router := mux.NewRouter()
	creds := []UserCredential{
		{
			Username: "user",
			Password: "valid_password",
			Roles:    []string{"User"},
		},
	}
	secret := "secret"
	api = New(router, creds, secret)
	api.Endpoints()
	os.Exit(m.Run())
}

func TestAPI_auth(t *testing.T) {
	data := authData{
		Username: "user",
		Password: "valid_password",
	}
	payload, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	got := rr.Body.String()
	want := "Token:"
	if !strings.HasPrefix(got, want) {
		t.Errorf("содержимое некорректно: получили %s, а хотели %s", got, want)
	}

	data = authData{
		Username: "user",
		Password: "wrong_password",
	}
	payload, _ = json.Marshal(data)
	req = httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(payload))
	rr = httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusUnauthorized)
	}
}
