package index

import (
	"gosearch/pkg/crawler"
	"reflect"
	"testing"
)

func TestStorage_getDoc(t *testing.T) {
	var docs = []crawler.Document{
		{
			ID:    1,
			URL:   "https://yandex.ru",
			Title: "Яндекс",
		},
		{
			ID:    2,
			URL:   "https://google.ru",
			Title: "Google",
		},
	}

	var s = New()
	for i := range docs {
		s.DocsTree.Insert(&docs[i])
	}

	var want = docs[1]
	got, err := s.getDoc(2)

	if !reflect.DeepEqual(want, *got) {
		t.Fatal("invalid result for existing doc")
	}

	got, err = s.getDoc(3)
	if (got != nil) || (err == nil) {
		t.Fatal("invalid returns for error")
	}
}

func TestStorage_addNewWords(t *testing.T) {
	var doc = crawler.Document{
		ID:    1,
		URL:   "https://site.com",
		Title: "Слово1 Слово2 Слово1 СЛОВО3",
	}

	var s = New()
	s.DocsTree.Insert(&doc)

	want := map[string][]int{
		"cлово1": []int{1},
		"cлово2": []int{1},
		"cлово3": []int{1},
	}

	s.addNewWords(doc)
	got := s.Index

	if len(want) != len(got) {
		t.Fatal("invalid word splitting")
	}
}
