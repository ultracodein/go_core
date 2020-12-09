package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// API реализует функционал аутентификации по JWT
type API struct {
	router *mux.Router
	creds  []UserCredential
	secret string
}

// UserCredential хранит аутентификационную информацию и роли пользователя
type UserCredential struct {
	Username string
	Password string
	Roles    []string
}

// authData задает формат запроса на аутентификацию
type authData struct {
	Username string
	Password string
}

// New - конструктор API
func New(router *mux.Router, creds []UserCredential, secret string) *API {
	api := API{
		router: router,
		creds:  creds,
		secret: secret,
	}
	return &api
}

// Endpoints регистрирует конечные точки API
func (api *API) Endpoints() {
	api.router.Handle("/auth", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(api.auth)))
}

func (api *API) auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "unsupported method", http.StatusInternalServerError)
		return
	}
	var data authData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roles, err := api.getUserRole(data.Username, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := api.getJWT(data.Username, roles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Token: " + token))
}

func (api *API) getUserRole(usr, pwd string) ([]string, error) {
	for _, cred := range api.creds {
		if cred.Username == usr && cred.Password == pwd {
			return cred.Roles, nil
		}
	}
	return nil, errors.New("ошибка авторизации")

}

func (api *API) getJWT(usr string, roles []string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"usr":   usr,
			"nbf":   time.Now().Unix(),
			"roles": roles,
		})

	return token.SignedString([]byte(api.secret))
}
