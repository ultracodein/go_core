package rpcsrv

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/engine"
)

// Service предоставляет интерфейс для удаленного вызова метода поиска документов
type Service struct {
	engine *engine.Service
}

// New - конструктор Service
func New(engine *engine.Service) *Service {
	srv := Service{
		engine: engine,
	}
	return &srv
}

// Search выполняет поиск документов
func (srv *Service) Search(query string, result *[]crawler.Document) error {
	docs := srv.engine.Search(query)
	*result = docs
	return nil
}
