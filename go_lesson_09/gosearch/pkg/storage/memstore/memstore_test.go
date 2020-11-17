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

func docs(count int) []crawler.Document {
	docs := make([]crawler.Document, count)
	for i := range docs {
		docs[i].ID = i
	}
	return docs
}

func benchmarkDBStoreDocs(count int, b *testing.B) {
	db := New()
	var docs = docs(count)

	for n := 0; n < b.N; n++ {
		db.StoreDocs(docs)
		for i := range docs {
			docs[i].ID++
		}
	}
}

func BenchmarkDB_StoreDocs1(b *testing.B)   { benchmarkDBStoreDocs(1, b) }
func BenchmarkDB_StoreDocs10(b *testing.B)  { benchmarkDBStoreDocs(10, b) }
func BenchmarkDB_StoreDocs100(b *testing.B) { benchmarkDBStoreDocs(100, b) }

var result []crawler.Document

func benchmarkDBDocs(count int, b *testing.B) {
	db := New()
	var docs = docs(count)
	db.StoreDocs(docs)
	var found []crawler.Document

	for n := 0; n < b.N; n++ {
		found = db.Docs([]int{0, count - 1})
	}

	result = found
}

func BenchmarkDB_Docs1(b *testing.B)   { benchmarkDBDocs(1, b) }
func BenchmarkDB_Docs10(b *testing.B)  { benchmarkDBDocs(10, b) }
func BenchmarkDB_Docs100(b *testing.B) { benchmarkDBDocs(100, b) }
