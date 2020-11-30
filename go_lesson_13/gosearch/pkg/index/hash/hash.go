package hash

import (
	"gosearch/pkg/crawler"
	"regexp"
	"strings"
)

// Index - индекс на основе хэш-таблицы.
type Index struct {
	Data map[string][]int
}

// New - конструктор.
func New() *Index {
	var ind Index
	ind.Data = make(map[string][]int)
	return &ind
}

// Add добавляет данные из переданных документов в индекс.
//
// Сначала происходит выделение лексем как ключей словаря из данных документа.
// Потом проверяется наличие номера документа в значении словаря для лексемы.
// Если номер документа не найден, то он добавляется в значение словаря.
func (index *Index) Add(docs []crawler.Document) {
	for _, doc := range docs {
		foundTokens := tokens(doc.Title)
		if foundTokens == nil {
			continue
		}

		for _, token := range foundTokens {
			if !exists(index.Data[token], doc.ID) {
				index.Data[token] = append(index.Data[token], doc.ID)
			}
		}
	}
}

// Search возвращает номера документов, где встречается данная лексема.
func (index *Index) Search(token string) []int {
	return index.Data[strings.ToLower(token)]
}

// Разделение строки на лексемы.
func tokens(s string) []string {
	wordPattern := regexp.MustCompile(`[\p{L}\d]+`)
	words := wordPattern.FindAllString(s, -1)
	if words == nil {
		return nil
	}

	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return words
}

// Проверка наличия элемента в массиве.
func exists(ids []int, item int) bool {
	for _, id := range ids {
		if id == item {
			return true
		}
	}
	return false
}
