package storage

// Хранилище отсканированных документов.

import (
	"gosearch/pkg/crawler"
)

// Interface определяет контракт хранилища данных.
type Interface interface {
	Docs([]int) []crawler.Document
	StoreDocs([]crawler.Document) error
	Doc(int) *crawler.Document
	Content() *[]crawler.Document
	Del(int) bool
}
