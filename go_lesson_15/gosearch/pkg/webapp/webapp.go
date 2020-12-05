package webapp

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gosearch/pkg/index"
	"gosearch/pkg/storage"
)

// Service - служба web-приложения
type Service struct {
	addr    string
	srv     http.Server
	index   *index.Interface
	storage *storage.Interface
}

// New - конструктор службы web-приложения
func New(addr string, index *index.Interface, storage *storage.Interface) *Service {
	s := Service{}
	s.addr = addr
	s.srv = http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      nil,
		Addr:         addr,
	}
	s.index = index
	s.storage = storage

	return &s
}

// Start запускает веб-службу
func (s *Service) Start() (*sync.WaitGroup, error) {
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		return nil, err
	}

	http.HandleFunc("/index", s.viewIndexHandler)
	http.HandleFunc("/docs", s.viewDocsHandler)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.srv.Serve(listener)
	}(&wg)

	return &wg, nil
}

func (s *Service) viewIndexHandler(w http.ResponseWriter, r *http.Request) {
	page := ""
	for word, ids := range (*s.index).Content() {
		docIds := ""
		for _, id := range ids {
			docIds += " " + strconv.Itoa(id)
		}

		page += fmt.Sprintf(`<li>%s:%s</li>`, word, docIds)
	}
	page = fmt.Sprintf(`<html><body><h2>Index</h2><ul>%s</ul></body></html>`, page)
	w.Write([]byte(page))
}

func (s *Service) viewDocsHandler(w http.ResponseWriter, r *http.Request) {
	page := ""
	for _, doc := range (*s.storage).Content() {
		page += fmt.Sprintf(`<li>%d: <a href="%s">%s</a></li>`, doc.ID, doc.URL, doc.Title)
	}
	page = fmt.Sprintf(`<html><body><h2>Docs</h2><ul>%s</ul></body></html>`, page)
	w.Write([]byte(page))
}
