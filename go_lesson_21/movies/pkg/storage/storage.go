package storage

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
	MovieBulkAdd(movies []Movie) ([]int, error)
	MovieDelete(movieID int) error
	MovieUpdate(movie Movie) error
	MovieGetAll(studioID int) ([]Movie, error)
}
