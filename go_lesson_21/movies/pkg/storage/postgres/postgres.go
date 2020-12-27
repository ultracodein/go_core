package postgres

import (
	"context"
	"movies/pkg/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB - хранилище данных
type DB struct {
	conn string
	pool *pgxpool.Pool
}

// New - конструктор
func New(conn string) *DB {
	db := DB{
		conn: conn,
		pool: nil,
	}
	return &db
}

// Connect выполняет подключение к БД
func (db *DB) Connect() error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, db.conn)
	if err != nil {
		return err
	}

	db.pool = pool
	return nil
}

// Disconnect закрывает соединения с БД
func (db *DB) Disconnect() {
	db.pool.Close()
}

// MovieBulkAdd сохраняет в БД массив фильмов
func (db *DB) MovieBulkAdd(ctx context.Context, movies []storage.Movie) ([]int, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)
	for _, movie := range movies {
		batch.Queue(
			`INSERT INTO movies(title, release_year, studio_id, gross, rating) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			movie.Title,
			movie.ReleaseYear,
			movie.StudioID,
			movie.Gross,
			movie.Rating,
		)
	}
	res := tx.SendBatch(ctx, batch)

	var ids []int
	for range movies {
		var id int
		err = res.QueryRow().Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	err = res.Close()
	if err != nil {
		return nil, err
	}

	return ids, tx.Commit(ctx)
}

// MovieDelete удаляет фильм
func (db *DB) MovieDelete(ctx context.Context, movieID int) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM movies WHERE id = $1`, movieID)
	return err
}

// MovieUpdate обновляет фильм
func (db *DB) MovieUpdate(ctx context.Context, movie storage.Movie) error {
	_, err := db.pool.Exec(
		ctx,
		`UPDATE movies
		SET title = $2,
		release_year = $3,
		studio_id = $4,
		gross = $5,
		rating = $6
		WHERE id = $1`,
		movie.ID,
		movie.Title,
		movie.ReleaseYear,
		movie.StudioID,
		movie.Gross,
		movie.Rating,
	)
	return err
}

// MovieGetAll возвращает все фильмы студии
func (db *DB) MovieGetAll(ctx context.Context, studioID int) ([]storage.Movie, error) {
	rows, err := db.pool.Query(
		ctx,
		`SELECT * FROM movies WHERE studio_id = $1
		OR studio_id IN (
			SELECT DISTINCT studio_id FROM movies
			WHERE studio_id <> $1 AND (SELECT $1) = 0
		)`,
		studioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []storage.Movie
	for rows.Next() {
		var m storage.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.StudioID, &m.Gross, &m.Rating)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return movies, nil
}
