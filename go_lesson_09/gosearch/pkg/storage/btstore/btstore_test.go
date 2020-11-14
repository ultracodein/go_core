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

func BenchmarkDB_StoreDocs(b *testing.B) {
	var docs = []crawler.Document{{ID: 1}, {ID: 2}}
	db := New()

	for n := 0; n < b.N; n++ {
		db.StoreDocs(docs)
		for i := range docs {
			docs[i].ID++
		}
	}
}

var result []crawler.Document

func BenchmarkDB_Docs(b *testing.B) {
	var docs = []crawler.Document{{ID: 1}, {ID: 2}, {ID: 3}}
	db := New()
	db.StoreDocs(docs)
	var found []crawler.Document

	for n := 0; n < b.N; n++ {
		found = db.Docs([]int{1, 3})
	}

	result = found
}
