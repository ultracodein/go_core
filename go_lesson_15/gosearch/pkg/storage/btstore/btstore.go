package btstore

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/storage/btstore/btree"
)

// DB - это хранилище данных на основе бинарного дерева
type DB struct {
	Tree      btree.Structure
	Documents []crawler.Document
}

// New - это конструктор хранилища данных
func New() *DB {
	db := DB{
		Tree:      btree.Structure{},
		Documents: make([]crawler.Document, 0),
	}
	return &db
}

// StoreDocs добавляет новые документы.
func (db *DB) StoreDocs(docs []crawler.Document) error {
	db.Documents = append(db.Documents, docs...)
	for i := range docs {
		db.Tree.Insert(&docs[i])
	}
	return nil
}

// Docs возвращает документы по их номерам.
func (db *DB) Docs(ids []int) []crawler.Document {
	var result []crawler.Document
	for _, id := range ids {
		vtx := db.Tree.Find(id)
		if vtx == nil {
			continue
		}
		item := *vtx.Item
		doc := item.(*crawler.Document)
		result = append(result, *doc)
	}
	return result
}

// Content возвращает содержимое хранилища
func (db *DB) Content() []crawler.Document {
	return db.Documents
}
