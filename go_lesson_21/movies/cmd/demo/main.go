package main

import (
	"context"
	"fmt"
	"log"
	"movies/pkg/storage"
	"movies/pkg/storage/postgres"
)

func main() {
	ctx := context.Background()
	db := postgres.New(ctx, "postgres://postgres:@localhost")

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Disconnect()

	movie := storage.Movie{
		ID:          0,
		Title:       "Бегущий по лезвию 2049",
		ReleaseYear: 2017,
		StudioID:    3,
		Gross:       259239658,
		Rating:      "PG-18",
	}

	// BulkAdd
	movies := []storage.Movie{
		movie,
	}
	ids, err := db.MovieBulkAdd(movies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ids)

	// GetAll
	movies, err = db.MovieGetAll(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", movies)

	// Update
	movie.ID = 4
	movie.Title = "Blade Runner 2049"
	err = db.MovieUpdate(movie)
	if err != nil {
		log.Fatal(err)
	}
	movies, err = db.MovieGetAll(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", movies)

	// Delete
	err = db.MovieDelete(4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted.\n")
}
