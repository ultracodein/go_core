package crawler

// Поисковый робот.
// Осуществляет сканирование сайтов.

// Interface определяет контракт поискового робота.
type Interface interface {
	Scan(urls []string, depth int) ([]Document, map[string]error)
}

// Document - документ, веб-страница, полученная поисковым роботом.
type Document struct {
	ID    int
	URL   string
	Title string
	Body  string
}

// Value возвращает ID документа
func (d *Document) Value() int {
	return d.ID
}
