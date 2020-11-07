package bststore

import (
	"fmt"
	"gosearch/pkg/crawler"
)

// DB - это хранилище данных на основе бинарного дерева
type DB struct {
	Tree      Tree
	Documents []crawler.Document
}

// Tree - это бинарное дерево
type Tree struct {
	Root *Vertex
	Len  int
}

// Vertex - это вершина графа, хранящая:
//  - "значение" хранимой в дереве сущности (Value)
//  - указатель на хранимую сущность (Item)
//  - указатели на вершины "левее" (L) и "правее" (R)
type Vertex struct {
	Value int
	Item  *Evaluator
	L     *Vertex
	R     *Vertex
}

// Evaluator абстрагирует сущность, которую необходимо поместить в дерево.
// Value() позволяет получить "значение" сущности (для возможности сравнения с другими вершинами).
type Evaluator interface {
	Value() int
}

// New - это конструктор хранилища данных
func New() *DB {
	db := DB{
		Tree:      Tree{},
		Documents: make([]crawler.Document, 0),
	}
	return &db
}

// StoreDocs добавляет новые документы.
func (db *DB) StoreDocs(docs []crawler.Document) error {
	if docs == nil {
		return nil
	}

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
		vtx, err := db.Tree.Find(id)
		if err != nil {
			continue
		}
		item := *vtx.Item
		doc := item.(*crawler.Document)
		result = append(result, *doc)
	}
	return result
}

// Insert реализует вставку сущности в дерево
func (t *Tree) Insert(e Evaluator) {
	// инициализируем новую (добавляемую) вершину
	var newValue = e.Value()
	var ins Vertex
	ins.Value = newValue
	ins.Item = &e

	// если в дереве нет вершин, добавляем новую вершину в корень
	if t.Root == nil {
		t.Len++
		t.Root = &ins
		return
	}

	// если в дереве есть вершины, обходим дерево вглубь (начиная с корня)
	var next = t.Root

	for {
		// если в дереве уже есть вершина с таким же значением - новую не добавляем
		if next.Value == newValue {
			break
		}

		var moveLeft = newValue < next.Value
		var moveRight = !moveLeft

		if moveLeft {
			if next.L == nil {
				next.L = &ins
				t.Len++
				break
			} else {
				next = next.L
				continue
			}
		}

		if moveRight {
			if next.R == nil {
				next.R = &ins
				t.Len++
				break
			} else {
				next = next.R
				continue
			}
		}
	}
}

// Find осуществляет поиск вершины по указанному значению
func (t *Tree) Find(value int) (*Vertex, error) {
	if t.Root == nil {
		return t.Root, fmt.Errorf("поиск невозможен - в дереве отсутствуют вершины")
	}

	var next = t.Root

	for {
		if next.Value == value {
			return next, nil
		}

		var moveLeft = value < next.Value
		var moveRight = !moveLeft

		if moveLeft {
			if next.L == nil {
				return next, fmt.Errorf("значение не найдено")
			}
			next = next.L
			continue
		}

		if moveRight {
			if next.R == nil {
				return next, fmt.Errorf("значение не найдено")
			}
			next = next.R
			continue
		}
	}
}
