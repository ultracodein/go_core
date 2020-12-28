package storage

import "context"

// Movie - кинофильм
type Movie struct {
	ID          int
	Title       string
	ReleaseYear int
	StudioID    int
	Gross       int
	Rating      string
}

// Interface определяет контракт хранилища данных.
type Interface interface {
	MovieBulkAdd(ctx context.Context, movies []Movie) ([]int, error)
	MovieDelete(ctx context.Context, movieID int) error
	MovieUpdate(ctx context.Context, movie Movie) error
	MovieGetAll(ctx context.Context, studioID int) ([]Movie, error)
}
