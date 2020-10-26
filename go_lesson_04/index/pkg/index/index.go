// index реализует индексирование загруженных web-страниц
// и предоставляет функционал поиска по данном индексу

package index

import (
	"errors"
	"fmt"
	"gosearch/pkg/crawler"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Storage - структура для хранения массива документов и обратного индекса
// (словаря, где ключом является слово из Title, а значением - ID документа)
type Storage struct {
	Docs    []crawler.Document
	Reverse map[string][]int
}

// Find формирует отчет о результатах поиска указанного слова
// (содержит ID и Title документа, в которых было найдено указанное слово)
func (s *Storage) Find(word string) string {
	docs, found := s.Reverse[word]
	if !found {
		return "Nothing.\n"
	}

	var report string = ""
	for _, id := range docs {
		doc, err := s.getDocBin(id)
		if err != nil {
			continue
		}

		report += "[" + strconv.Itoa(doc.ID) + "] " + doc.Title + "\n"
	}

	return report
}

// Fill загружает содержимое указанных web-страниц и индексирует полученные данные
func Fill(scn crawler.Scanner, urls []string, depth int) *Storage {
	var storage = Storage{
		Docs:    make([]crawler.Document, 0),
		Reverse: make(map[string][]int),
	}

	var scanData = make([]crawler.Document, 0)

	for i, url := range urls {
		fmt.Printf("Scanning URL #%d: %s\n", i+1, url)

		data, err := scn.Scan(url, depth)
		if err != nil {
			continue
		}

		scanData = append(scanData, data...)
	}

	storage.append(scanData)

	return &storage
}

func (s *Storage) append(data []crawler.Document) {
	// инициализируем генератор случайных чисел
	seed := rand.NewSource(time.Now().UnixNano())
	rgen := rand.New(seed)

	// формируем массив случайных ID для новой порции документов
	var newCount = len(data)
	var existCount = len(s.Docs)
	var id int
	deltas := rgen.Perm(newCount)

	for i, doc := range data {
		id = existCount + 1 + deltas[i]

		// присваиваем документу случайный ID и добавляем его в массив
		doc.ID = id
		s.Docs = append(s.Docs, doc)

		// добавляем в обратный индекс новые слова
		s.addNewWords(doc)
	}

	// сортируем массив документов
	sort.Slice(s.Docs, func(i, j int) bool { return s.Docs[i].ID < s.Docs[j].ID })
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
		if _, duplicate := uniqueWords[word]; duplicate {
			continue
		}
		uniqueWords[word] = true

		// добавляем документ в список для данного слова
		s.Reverse[word] = append(s.Reverse[word], doc.ID)
	}
}

// getDocBin возвращает указатель на документ с указанным ID (бинарный поиск)
func (s *Storage) getDocBin(id int) (*crawler.Document, error) {
	i := sort.Search(len(s.Docs), func(i int) bool { return s.Docs[i].ID >= id })
	if i < len(s.Docs) && s.Docs[i].ID == id {
		return &s.Docs[i], nil
	}

	return nil, errors.New("not found")
}

// getDoc возвращает указатель на документ с указанным ID (простой перебор массива)
func (s *Storage) getDoc(id int) (*crawler.Document, error) {
	for _, doc := range s.Docs {
		if doc.ID == id {
			return &doc, nil
		}
	}

	return nil, errors.New("not found")
}
