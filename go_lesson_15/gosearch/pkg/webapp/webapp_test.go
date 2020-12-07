package webapp

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btstore"
	"gosearch/pkg/storage/btstore/btree"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var webapp *Service

func TestMain(m *testing.M) {
	// инициализируем индекс и хранилище тестовыми значениями
	var index index.Interface
	var storage storage.Interface
	index = &hash.Index{
		Data: map[string][]int{
			"word1": []int{0, 1},
			"word2": []int{2},
		},
	}
	storage = &btstore.DB{
		Tree: btree.Structure{},
		Documents: []crawler.Document{
			crawler.Document{ID: 0, URL: "url1", Title: "word1"},
			crawler.Document{ID: 1, URL: "url2", Title: "word1"},
			crawler.Document{ID: 2, URL: "url3", Title: "word2"},
		},
	}
	webapp = New(":8000", &index, &storage)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestService_viewIndexHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webapp.viewIndexHandler)
	handler.ServeHTTP(rr, req)

	gotStatus := rr.Code
	wantStatus := http.StatusOK
	if gotStatus != wantStatus {
		t.Errorf("код неверен: получили %d, а хотели %d", gotStatus, wantStatus)
	}

	got := rr.Body.String()
	wants := []string{
		"word1: 0 1",
		"word2: 2",
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("содержимое некорректно: получили %s, а хотели %s", got, want)
		}
	}
}

func TestService_viewDocsHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/docs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webapp.viewDocsHandler)
	handler.ServeHTTP(rr, req)

	gotStatus := rr.Code
	wantStatus := http.StatusOK
	if gotStatus != wantStatus {
		t.Errorf("код неверен: получили %d, а хотели %d", gotStatus, wantStatus)
	}

	got := rr.Body.String()
	wants := []string{
		`0: <a href="url1">word1</a>`,
		`1: <a href="url2">word1</a>`,
		`2: <a href="url3">word2</a>`,
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("содержимое некорректно: получили %s, а хотели %s", got, want)
		}
	}
}
