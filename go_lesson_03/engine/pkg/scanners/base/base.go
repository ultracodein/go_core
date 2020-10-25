package base

import "crawler/pkg/spider"

// SiteScanner является сканером страниц по умолчанию
type SiteScanner func(string, int)

// Scan реализует интерфейс Scanner для сканера страниц по умолчанию
func (s SiteScanner) Scan(url string, depth int) (data map[string]string, err error) {
	return spider.Scan(url, depth)
}
