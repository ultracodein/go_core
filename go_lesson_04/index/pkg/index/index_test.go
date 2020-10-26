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
		Title: "Слово1 Слово2 Слово1 Слово3",
	}

	want := map[string][]int{
		"cлово1": []int{1},
		"cлово2": []int{1},
		"cлово3": []int{1},
	}

	store.addNewWords(doc)
	got := store.Reverse

	if len(want) != len(got) {
		t.Fatal("invalid word splitting")
	}

	for k := range want {
		gotIds, ok := want[k]

		if !ok {
			t.Fatal("different indexes")
		}

		for i := 0; i < len(want[k]); i++ {
			if want[k][i] != gotIds[i] {
				t.Fatal("found duplicates")
			}
		}
	}
}
