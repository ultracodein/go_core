module auth/cmd/auth

go 1.15

replace auth/pkg/api => ../../pkg/api

require (
	auth/pkg/api v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0
)
