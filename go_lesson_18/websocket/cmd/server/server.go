package main

import (
	"net/http"
	"websocket/pkg/api"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	router := mux.NewRouter()
	upgrader := websocket.Upgrader{
		//CheckOrigin: func(r *http.Request) bool {
		//	return true
		//},
	}
	secret := "password"

	api := api.New(router, upgrader, secret)
	api.Endpoints()

	http.ListenAndServe(":8000", router)
}
