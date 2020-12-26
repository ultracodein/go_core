package postgres

import (
	"context"
	"fmt"
	"movies/pkg/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB - хранилище данных
type DB struct {
	ctx  context.Context
	conn string
	pool *pgxpool.Pool
}

// New - конструктор
func New(ctx context.Context, conn string) *DB {
	db := DB{
		conn: conn,
		ctx:  ctx,
		pool: nil,
	}
	return &db
}

// Connect выполняет подключение к БД
func (db *DB) Connect() error {
	pool, err := pgxpool.Connect(db.ctx, db.conn)
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
func (db *DB) MovieBulkAdd(movies []storage.Movie) ([]int, error) {
	tx, err := db.pool.Begin(db.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(db.ctx)

	batch := new(pgx.Batch)
	batch.Queue(`SELECT SETVAL('movies_id_seq', (SELECT MAX(id) FROM movies))`)
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

	res := tx.SendBatch(db.ctx, batch)

	ids := make([]int, 0)
	for i := 1; i <= len(movies)+1; i++ {
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

	return ids[1:], tx.Commit(db.ctx)
}

// MovieDelete удаляет фильм
func (db *DB) MovieDelete(movieID int) error {
	_, err := db.pool.Exec(db.ctx, `DELETE FROM movies WHERE id = $1`, movieID)
	return err
}

// MovieUpdate обновляет фильм
func (db *DB) MovieUpdate(movie storage.Movie) error {
	_, err := db.pool.Exec(
		db.ctx,
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
func (db *DB) MovieGetAll(studioID int) ([]storage.Movie, error) {
	query := `SELECT id, title, release_year, studio_id, gross, rating FROM movies`
	if studioID > 0 {
		query += fmt.Sprintf(` WHERE studio_id = %d`, studioID)
	}

	rows, err := db.pool.Query(db.ctx, query)
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
