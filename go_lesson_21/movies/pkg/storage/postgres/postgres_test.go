package postgres

import (
	"context"
	"log"
	"movies/pkg/storage"
	"os"
	"reflect"
	"testing"
)

var (
	db       *DB
	ctx      context.Context
	newMovie storage.Movie
)

func TestMain(m *testing.M) {
	db = New("postgres://postgres:@localhost")
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	ctx = context.Background()
	newMovie = storage.Movie{
		ID:          0,
		Title:       "Бегущий по лезвию 2049",
		ReleaseYear: 2017,
		StudioID:    3,
		Gross:       259239658,
		Rating:      "PG-18",
	}

	os.Exit(m.Run())
}

func TestDB_MovieGetAll(t *testing.T) {
	tests := []struct {
		name     string
		studioID int
		want     int
	}{
		{
			name:     "All studios",
			studioID: 0,
			want:     3,
		},
		{
			name:     "Studio with id = 3",
			studioID: 3,
			want:     1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movies, _ := db.MovieGetAll(ctx, tt.studioID)
			got := len(movies)
			if got != tt.want {
				t.Errorf("получили список фильмов длиной %d, а хотели %d", got, tt.want)
				return
			}
		})
	}
}

func TestDB_MovieBulkAdd(t *testing.T) {
	movies := []storage.Movie{
		newMovie,
	}
	ids, _ := db.MovieBulkAdd(ctx, movies)
	newMovie.ID = ids[0]
	movies, _ = db.MovieGetAll(ctx, 3)

	got := movies[1]
	want := newMovie
	if !reflect.DeepEqual(got, want) {
		t.Errorf("был добавлен фильм %v, а хотели %v", got, want)
	}

	_ = db.MovieDelete(ctx, 4)
}

func TestDB_MovieUpdate(t *testing.T) {
	movies, _ := db.MovieGetAll(ctx, 2)
	movie := movies[0]
	want := "Diarios de motocicleta"
	movie.Title = want

	_ = db.MovieUpdate(ctx, movie)
	movies, _ = db.MovieGetAll(ctx, 2)
	movie = movies[0]
	got := movie.Title

	if got != want {
		t.Errorf("название изменено на %v, а хотели на %v", got, want)
	}

	movie.Title = "Че Гевара: Дневники мотоциклиста"
	_ = db.MovieUpdate(ctx, movie)
}

func TestDB_MovieDelete(t *testing.T) {
	movies := []storage.Movie{
		newMovie,
	}
	ids, _ := db.MovieBulkAdd(ctx, movies)
	id := ids[0]

	_ = db.MovieDelete(ctx, id)

	movies, _ = db.MovieGetAll(ctx, 0)
	got := len(movies)
	want := 3

	if got != want {
		t.Errorf("получили список фильмов длиной %d, а хотели %d", got, want)
	}
}
