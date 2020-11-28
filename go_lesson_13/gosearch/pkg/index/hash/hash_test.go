package hash

import (
	"gosearch/pkg/crawler"
	"reflect"
	"testing"
)

func TestIndex_Add(t *testing.T) {
	ind := New()
	docs := []crawler.Document{
		{
			ID:    10,
			Title: "Заголовок: заголовок",
		},
		{
			ID:    20,
			Title: "А это другой заголовок",
		},
		{
			ID:    30,
			Title: "",
		},
	}
	ind.Add(docs)

	// проверяем число ключей
	got := len(ind.Data)
	want := 4
	if got != want {
		t.Fatalf("получили %d, ожидалось %d", got, want)
	}

	// проверяем отсутствие дублей в значении
	got = len(ind.Data["заголовок"])
	want = 2
	if got != want {
		t.Fatalf("получили %d, ожидалось %d", got, want)
	}
}

func TestIndex_Search(t *testing.T) {
	ind := New()
	docs := []crawler.Document{
		{
			ID:    10,
			Title: "заголовок маленькими буквами",
		},
		{
			ID:    20,
			Title: "BIG ЗАГОЛОВОК",
		},
		{
			ID:    30,
			Title: "ЗаголовОК со знаками препинания!",
		},
	}
	ind.Add(docs)

	tests := []struct {
		name  string
		token string
		want  []int
	}{
		{
			name:  "Тест №1",
			token: "заголовок",
			want:  []int{10, 20, 30},
		},
		{
			name:  "Тест №2",
			token: "nil",
			want:  nil,
		},
		{
			name:  "Тест №3",
			token: "big",
			want:  []int{20},
		},
		{
			name:  "Тест №4",
			token: "препинания",
			want:  []int{30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ind.Search(tt.token); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("получили %v, ожидалось %v", got, tt.want)
			}
		})
	}
}
