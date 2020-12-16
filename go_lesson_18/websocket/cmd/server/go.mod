module websocket/cmd/server

go 1.15

replace websocket/pkg/api => ../../pkg/api

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	websocket/pkg/api v0.0.0-00010101000000-000000000000
)
