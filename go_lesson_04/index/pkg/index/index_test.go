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

	var store Storage
	store.Docs = docs

	var want = docs[1]
	got, err := store.getDoc(2)

	if !reflect.DeepEqual(want, *got) || err != nil {
		t.Fatal("getDoc returned invalid result")
	}

	got, err = store.getDocBin(2)
	if !reflect.DeepEqual(want, *got) || err != nil {
		t.Fatal("getDocBin returned invalid result")
	}
}

func TestStorage_addNewWords(t *testing.T) {
	var store Storage
	store.Reverse = make(map[string][]int)

	var doc = crawler.Document{
		ID:    1,
		URL:   "https://site.com",
		Title: "Слово1 Слово2 Слово1",
	}

	store.addNewWords(doc)
	got := store.Reverse

	if len(got) != 2 {
		t.Fatal("invalid word splitting")
	}
}
