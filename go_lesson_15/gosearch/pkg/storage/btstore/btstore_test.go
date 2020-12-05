package btstore

import (
	"gosearch/pkg/crawler"
	"reflect"
	"testing"
)

func TestDB_StoreDocs(t *testing.T) {
	var docs = []crawler.Document{{ID: 1}, {ID: 2}}
	db := New()

	db.StoreDocs(docs)
	got := db.Tree.Len
	want := 2

	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestDB_Docs(t *testing.T) {
	var docs = []crawler.Document{{ID: 1}, {ID: 2}, {ID: 3}}
	db := New()

	db.StoreDocs(docs)
	got := db.Docs([]int{1, 3, 777})
	want := []crawler.Document{{ID: 1}, {ID: 3}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
