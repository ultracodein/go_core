// index реализует индексирование загруженных web-страниц
// и предоставляет функционал поиска по данном индексу

package index

import (
	"fmt"
	"gosearch/pkg/bst"
	"gosearch/pkg/crawler"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Storage - структура для хранения дерева документов и обратного индекса
// (словаря, где ключом является слово из Title, а значением - ID документа)
type Storage struct {
	DocsTree bst.Tree
	Index    map[string][]int
}

// New - конструктор Storage
func New() *Storage {
	s := Storage{}
	s.Index = make(map[string][]int)
	return &s
}

// Fill загружает содержимое web-страниц и индексирует полученные данные
func (s *Storage) Fill(scn crawler.Scanner, urls []string, depth int) error {
	var scannedDocs = make([]crawler.Document, 0)

	for _, url := range urls {
		docs, err := scn.Scan(url, depth)
		if err != nil {
			continue
		}

		scannedDocs = append(scannedDocs, docs...)
	}

	if len(scannedDocs) == 0 {
		return fmt.Errorf("no new data")
	}

	s.addToIndex(scannedDocs)
	return nil
}

// Find возвращает указатели на документы, в которых было найдено искомое слово
func (s *Storage) Find(word string) []*crawler.Document {
	docs := s.Index[strings.ToLower(word)]

	if docs == nil {
		return nil
	}

	found := make([]*crawler.Document, 0)
	for _, id := range docs {
		doc, err := s.getDoc(id)
		if err != nil {
			continue
		}

		found = append(found, doc)
	}

	return found
}

func (s *Storage) addToIndex(docs []crawler.Document) {
	// инициализируем генератор случайных чисел
	seed := rand.NewSource(time.Now().UnixNano())
	rgen := rand.New(seed)

	// формируем массив случайных ID для новой порции документов
	var newCount = len(docs)
	var existCount = s.DocsTree.Len
	var id int
	deltas := rgen.Perm(newCount)

	for i := range docs {
		id = existCount + 1 + deltas[i]

		// присваиваем документу случайный ID и добавляем его в массив
		docs[i].ID = id
		s.DocsTree.Insert(&docs[i])

		// добавляем в обратный индекс новые слова
		s.addNewWords(docs[i])
	}
}

// addNewWords добавляет в словарь новые слова в нижнем регистре
func (s *Storage) addNewWords(doc crawler.Document) {
	wordPattern := regexp.MustCompile("[\\p{L}\\d]+")
	var titleLow = strings.ToLower(doc.Title)

	words := wordPattern.FindAllString(titleLow, -1)
	if words == nil {
		return
	}

	uniqueWords := make(map[string]bool)

	for _, word := range words {
		// исключаем повторную обработку слов по тому же документу
		if uniqueWords[word] {
			continue
		}
		uniqueWords[word] = true

		// добавляем документ в список для данного слова
		s.Index[word] = append(s.Index[word], doc.ID)
	}
}

// getDoc осуществляет поиск по бинарному дереву и возвращает представление найденного документа
func (s *Storage) getDoc(id int) (*crawler.Document, error) {
	vertex, err := s.DocsTree.Find(id)

	if err != nil {
		return nil, err
	}

	i := *vertex.Item
	return i.(*crawler.Document), nil
}
