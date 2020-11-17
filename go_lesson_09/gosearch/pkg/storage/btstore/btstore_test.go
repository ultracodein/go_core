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
