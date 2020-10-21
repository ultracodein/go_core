package main

import (
	"reflect"
	"testing"
)

func Test_collectData(t *testing.T) {
	urls := []string{"https://habr.com/", "https://www.cnews.ru/"}
	var depth int = 1

	var want = map[string]string{
		"https://habr.com/":     "Лучшие публикации за сутки / Хабр",
		"https://www.cnews.ru/": "Интернет-издание о высоких технологиях - CNews",
	}

	var got = make(map[string]string)
	collectData(urls, depth, got)

	if !reflect.DeepEqual(want, got) {
		t.Fatal("storage is not equal to wanted")
	}
}

func Test_findRelatedPages(t *testing.T) {
	var storageExample = map[string]string{
		"1": "Новости IT-технологий",
		"2": "It works",
	}

	type args struct {
		storage map[string]string
		search  string
	}
	tests := []struct {
		name      string
		args      args
		wantFound []string
	}{
		{
			name:      "Test All Results (Case-Insensitive)",
			args:      args{storage: storageExample, search: "it"},
			wantFound: []string{"1", "2"},
		},
		{
			name:      "Test One Result",
			args:      args{storage: storageExample, search: "технолог"},
			wantFound: []string{"1"},
		},
		{
			name:      "Test No Results",
			args:      args{storage: storageExample, search: "Go"},
			wantFound: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFound := findRelatedPages(tt.args.storage, tt.args.search); !reflect.DeepEqual(gotFound, tt.wantFound) {
				t.Errorf("findRelatedPages() = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
