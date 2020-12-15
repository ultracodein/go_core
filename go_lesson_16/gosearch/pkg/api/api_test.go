package api

import (
	"bytes"
	"encoding/json"
	"gosearch/pkg/crawler"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/memstore"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var api *API

func TestMain(m *testing.M) {
	router := mux.NewRouter()
	var index index.Interface
	var storage storage.Interface
	index = &hash.Index{
		Data: map[string][]int{
			"word1": []int{0},
			"word2": []int{1, 2},
		},
	}
	storage = memstore.New()
	docs := []crawler.Document{
		crawler.Document{ID: 0, URL: "url1", Title: "word1"},
		crawler.Document{ID: 1, URL: "url2", Title: "word2"},
		crawler.Document{ID: 2, URL: "url3", Title: "word2"},
	}
	storage.StoreDocs(docs)
	engine := engine.New(index, storage)

	api = New(router, engine, &index, &storage)
	api.Endpoints()
	os.Exit(m.Run())
}

func TestAPI_search(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/word2", nil)
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	wantIds := []int{1, 2}
	wantBytes, _ := json.Marshal(wantIds)
	want := string(wantBytes)
	got := strings.TrimSuffix(rr.Body.String(), "\n")

	if got != want {
		t.Errorf("содержимое некорректно: получили %s а хотели %s", got, want)
	}
}

func TestAPI_viewDocs(t *testing.T) {
	ids := []int{0, 2}
	payload, _ := json.Marshal(ids)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/docs/view", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	wantDocs := []crawler.Document{
		{
			ID:    0,
			URL:   "url1",
			Title: "word1",
			Body:  "",
		},
		{
			ID:    2,
			URL:   "url3",
			Title: "word2",
			Body:  "",
		},
	}
	wantBytes, _ := json.Marshal(wantDocs)
	want := string(wantBytes)
	got := strings.TrimSuffix(rr.Body.String(), "\n")

	if got != want {
		t.Errorf("содержимое некорректно: получили %s а хотели %s", got, want)
	}
}

func TestAPI_createDocs(t *testing.T) {
	doc := crawler.Document{
		ID:    999,
		URL:   "url9",
		Title: "word9",
		Body:  "",
	}
	payload, _ := json.Marshal(doc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/docs/create", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	want := crawler.Document{
		ID:    3,
		URL:   "url9",
		Title: "word9",
		Body:  "",
	}
	got := (*api.storage).Doc(3)

	if got == nil || !reflect.DeepEqual(*got, want) {
		t.Errorf("документ не добавлен или добавлен некорректно: получили %v а хотели %v", got, want)
	}
}

func TestAPI_updateDocs(t *testing.T) {
	origin := (*api.storage).Doc(1)
	doc := crawler.Document{
		ID:    1,
		URL:   "url2",
		Title: "word2",
		Body:  "body",
	}
	payload, _ := json.Marshal(doc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/docs/update/1", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	want := "body"
	got := origin.Body

	if got != want {
		t.Errorf("содержимое некорректно: получили %s а хотели %s", got, want)
	}
}

func TestAPI_deleteDocs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs/delete/1", nil)
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	got := (*api.storage).Doc(1)

	if got != nil {
		t.Errorf("содержимое некорректно: получили %v а хотели %v", got, nil)
	}
}
