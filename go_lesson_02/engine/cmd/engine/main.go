package main

import (
	"crawler/pkg/spider"
	"flag"
	"fmt"
	"regexp"
)

func main() {
	urls := []string{"https://habr.com/", "https://www.cnews.ru/"}
	var depth int = 1

	search := parseSearchPhrase()

	fmt.Printf("Today we have %d sites to scan!\n", len(urls))
	storage := make(map[string]string)
	collectData(urls, depth, storage)

	found := findRelatedPages(storage, search)

	fmt.Println("Pages found:")

	if len(found) > 0 {
		for _, link := range found {
			fmt.Printf("[%s] %s\n", storage[link], link)
		}
	} else {
		fmt.Println("Nothing.")
	}
}

func parseSearchPhrase() string {
	searchPtr := flag.String("search", "", "word/phrase to search in page title")
	flag.Parse()

	if *searchPtr != "" {
		fmt.Printf("Your search phrase is: %s\n", *searchPtr)
	}

	return *searchPtr
}

func collectData(urls []string, depth int, storage map[string]string) {
	for i, url := range urls {
		fmt.Printf("Scanning URL #%d: %s\n", i+1, url)

		data, err := spider.Scan(url, depth)

		if err != nil {
			continue
		}

		for link, title := range data {
			storage[link] = title
		}
	}
}

func findRelatedPages(storage map[string]string, search string) (found []string) {
	found = []string{}
	escapedSearch := "(?i)" + regexp.QuoteMeta(search)
	searchRegEx, err := regexp.Compile(escapedSearch)

	if err != nil {
		return
	}

	for link, title := range storage {
		if searchRegEx.Match([]byte(title)) {
			found = append(found, link)
		}
	}

	return
}
