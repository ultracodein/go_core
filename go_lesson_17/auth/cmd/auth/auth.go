package main

import (
	"auth/pkg/api"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	creds := initCreds()
	secret := "auth-service-password"
	api := api.New(router, creds, secret)
	api.Endpoints()

	http.ListenAndServe(":8000", router)
}

func initCreds() []api.UserCredential {
	var c = []api.UserCredential{
		{
			Username: "user",
			Password: "user",
			Roles:    []string{"User"},
		},
		{
			Username: "oper",
			Password: "oper",
			Roles:    []string{"Manager", "Backup"},
		},
		{
			Username: "admin",
			Password: "admin",
			Roles:    []string{"Backup", "Admin"},
		},
	}
	return c
}
