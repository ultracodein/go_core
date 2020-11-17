package engine

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Engine - поисковый движок.
// Его задача - обслуживание поисковых запросов.
// функциональность:
// - обработка поискового запроса;
// - поиск документов в индексе;
// - запрос документов из хранилища;
// - возврат посиковой выдачи.

// Service - поисковый движок.
type Service struct {
	index   index.Interface
	storage storage.Interface
}

// engine_query_len_sum/engine_query_len_count - средняя длина за все время
// rate(engine_query_len_sum[5m])/rate(engine_query_len_count[5m]) - средняя длина за последние 5 минут
var averageQueryLen = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "engine_query_len",
	Help:    "Длина поискового запроса, Байт.",
	Buckets: prometheus.LinearBuckets(0, 1, 50),
})

// New - конструктор.
func New(index index.Interface, storage storage.Interface) *Service {
	s := Service{
		index:   index,
		storage: storage,
	}
	return &s
}

// Search ищет документы, соответствующие поисковому запросу.
func (s *Service) Search(query string) []crawler.Document {
	if query == "" {
		return nil
	}
	averageQueryLen.Observe((float64(len([]byte(query)))))
	ids := s.index.Search(query)
	docs := s.storage.Docs(ids)
	return docs
}
