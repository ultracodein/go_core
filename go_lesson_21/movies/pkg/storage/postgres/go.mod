module movies/pkg/storage/postgres

replace movies/pkg/storage => ../

go 1.15

require (
	github.com/jackc/pgx/v4 v4.10.1
	movies/pkg/storage v0.0.0-00010101000000-000000000000
)
