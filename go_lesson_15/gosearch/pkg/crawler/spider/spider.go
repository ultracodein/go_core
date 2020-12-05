// Package spider реализует сканер содержимого веб-сайтов.
// Пакет позволяет получить список ссылок и заголовков страниц внутри веб-сайта по его URL.
package spider

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"gosearch/pkg/crawler"

	"golang.org/x/net/html"
)

// Service - служба поискового робота.
type Service struct {
	maxThreads int
}

type result struct {
	url  string
	docs []crawler.Document
	err  error
}

// New - конструктор службы поискового робота.
func New(maxThreads int) *Service {
	s := Service{}
	s.maxThreads = maxThreads
	return &s
}

// Scan осуществляет обход сайтов в многопоточном режиме
func (s *Service) Scan(urls []string, depth int) (docs []crawler.Document, errors map[string]error) {
	var wg sync.WaitGroup
	in := make(chan string)
	out := make(chan result)

	// заполняем входной канал url
	go func() {
		for _, url := range urls {
			in <- url
		}
	}()

	// определяем число потоков и запускаем их
	urlCount := len(urls)
	count := s.getThreadCount(urlCount)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go s.scanThread(in, out, depth, &wg)
	}

	// читаем и обрабатываем результаты из выходного канала
	errors = make(map[string]error)
	for range urls {
		res := <-out

		if res.err != nil {
			errors[res.url] = res.err
			continue
		}

		docs = append(docs, res.docs...)
	}

	close(in)
	wg.Wait()

	return docs, errors
}

func (s *Service) getThreadCount(urlCount int) int {
	threadCount := 0
	if urlCount > s.maxThreads {
		threadCount = s.maxThreads
	} else {
		threadCount = urlCount
	}
	return threadCount
}

func (s *Service) scanThread(in <-chan string, out chan<- result, depth int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		url, open := <-in
		if !open {
			return
		}

		log.Printf("Идет сканирование %s\n", url)
		docs, err := scanPage(url, depth)
		out <- result{
			url:  url,
			docs: docs,
			err:  err,
		}
	}
}

func scanPage(url string, depth int) (docs []crawler.Document, err error) {
	pages := make(map[string]string)

	err = parse(url, url, depth, pages)
	if err != nil {
		return nil, err
	}

	for url, title := range pages {
		doc := crawler.Document{
			URL:   url,
			Title: title,
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// parse рекурсивно обходит ссылки на странице, переданной в url.
// Глубина рекурсии задаётся в depth.
// Каждая найденная ссылка записывается в ассоциативный массив
// data вместе с названием страницы.
func parse(url, baseurl string, depth int, data map[string]string) error {
	if depth == 0 {
		return nil
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	page, err := html.Parse(response.Body)
	if err != nil {
		return err
	}

	data[url] = pageTitle(page)

	if depth == 1 {
		return nil
	}
	links := pageLinks(nil, page)
	for _, link := range links {
		link = strings.TrimSuffix(link, "/")
		// относительная ссылка
		if strings.HasPrefix(link, "/") && len(link) > 1 {
			link = baseurl + link
		}
		// ссылка уже отсканирована
		if data[link] != "" {
			continue
		}
		// ссылка содержит базовый url полностью
		if strings.HasPrefix(link, baseurl) {
			parse(link, baseurl, depth-1, data)
		}
	}

	return nil
}

// pageTitle осуществляет рекурсивный обход HTML-страницы и возвращает значение элемента <tittle>.
func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

// pageLinks рекурсивно сканирует узлы HTML-страницы и возвращает все найденные ссылки без дубликатов.
func pageLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !SliceContains(links, a.Val) {
					links = append(links, a.Val)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = pageLinks(links, c)
	}
	return links
}

// SliceContains возвращает true если массив содержит переданное значение
func SliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
