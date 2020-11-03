package bststore

import (
	"gosearch/pkg/crawler"
	"testing"
)

func genTree() *Tree {
	var db = New()
	tree := db.Tree

	var docs = []crawler.Document{
		{ID: 8}, {ID: 8},
		{ID: 3}, {ID: 10},
		{ID: 1}, {ID: 6}, {ID: 14},
		{ID: 4}, {ID: 7}, {ID: 13},
	}

	for _, doc := range docs {
		tree.Insert(&doc)
	}

	return &tree
}

func TestTree_Insert(t *testing.T) {
	tree := genTree()

	b1 := tree.Root.L.L.Value
	b2 := tree.Root.L.R.L.Value
	b3 := tree.Root.L.R.R.Value
	b4 := tree.Root.R.R.L.Value

	var want = [4]int{1, 4, 7, 13}
	var got = [4]int{b1, b2, b3, b4}

	if got != want {
		t.Fatalf("want %v, got %v", want, got)
	}

	if tree.Len != 9 {
		t.Fatalf("want 9, got %v", tree.Len)
	}
}

func TestTree_Find(t *testing.T) {
	tree := genTree()

	_, err := tree.Find(0)
	if err == nil {
		t.Fatalf("want error, got nil")
	}

	vertex, err := tree.Find(4)
	if err != nil {
		t.Fatalf("want nil, got error")
	}

	if vertex.Value != 4 {
		t.Fatalf("want vertex = 4, got vertex = %d", vertex.Value)
	}
}

func TestDB_Docs(t *testing.T) {
	var docs = []crawler.Document{{ID: 1}, {ID: 2}}
	want := len(docs)

	db := New()
	db.StoreDocs(docs)
	got := len(db.Docs([]int{1, 2}))

	if got != want {
		t.Fatalf("want %d, got %d", want, got)
	}
}
