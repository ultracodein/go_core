package memstore

import (
	"gosearch/pkg/crawler"
	"testing"
)

func TestDB_StoreDocs(t *testing.T) {
	db := New()
	docs := []crawler.Document{{ID: 10}, {ID: 20}}
	db.StoreDocs(docs)
	got := len(db.docs)
	want := 2
	if got != want {
		t.Fatalf("получили %d, ожидалось %d", got, want)
	}
}

func TestDB_Docs(t *testing.T) {
	db := New()
	docs := []crawler.Document{{ID: 10}, {ID: 20}}
	db.StoreDocs(docs)
	got := len(db.Docs([]int{10, 20}))
	want := 2
	if got != want {
		t.Fatalf("получили %d, ожидалось %d", got, want)
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
