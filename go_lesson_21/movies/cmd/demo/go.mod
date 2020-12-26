module movies/cmd/demo

replace (
	movies/pkg/storage => ../../pkg/storage
	movies/pkg/storage/postgres => ../../pkg/storage/postgres
)

go 1.15

require (
	github.com/jackc/pgx/v4 v4.10.1 // indirect
	movies/pkg/storage v0.0.0-00010101000000-000000000000
	movies/pkg/storage/postgres v0.0.0-00010101000000-000000000000
)
