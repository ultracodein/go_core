package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gosearch/pkg/crawler"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"

	"github.com/gorilla/mux"
)

// API предоставляет программный интерфейс для работы с индексом и хранилищем
type API struct {
	router  *mux.Router
	engine  *engine.Service
	index   *index.Interface
	storage *storage.Interface
}

// New - конструктор API
func New(router *mux.Router, engine *engine.Service, index *index.Interface, storage *storage.Interface) *API {
	api := API{
		router:  router,
		engine:  engine,
		index:   index,
		storage: storage,
	}
	return &api
}

// Endpoints регистрирует конечные точки API
func (api *API) Endpoints() {
	api.router.HandleFunc("/api/v1/search/{word}", api.search).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/api/v1/docs/view", api.viewDocs).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/api/v1/docs/create", api.createDocs).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/api/v1/docs/update/{id:[0-9]+}", api.updateDocs).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/api/v1/docs/delete/{id:[0-9]+}", api.deleteDocs).Methods(http.MethodGet, http.MethodOptions)
}

func (api *API) search(w http.ResponseWriter, r *http.Request) {
	word := mux.Vars(r)["word"]
	docs := api.engine.Search(word)
	var ids []int
	for _, doc := range docs {
		ids = append(ids, doc.ID)
	}

	err := json.NewEncoder(w).Encode(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) viewDocs(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	docs := (*api.storage).Docs(ids)
	err = json.NewEncoder(w).Encode(docs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) createDocs(w http.ResponseWriter, r *http.Request) {
	var doc crawler.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	docs := *(*api.storage).Content()
	doc.ID = len(docs)
	err = (*api.storage).StoreDocs([]crawler.Document{doc})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) updateDocs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var doc crawler.Document
	err = json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	origin := (*api.storage).Doc(id)
	if origin == nil {
		http.Error(w, "документ не найден", http.StatusInternalServerError)
		return
	}

	*origin = doc
}

func (api *API) deleteDocs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deleted := (*api.storage).Del(id)
	if !deleted {
		http.Error(w, "документ не найден", http.StatusInternalServerError)
		return
	}
}
