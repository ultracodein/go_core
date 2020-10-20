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

	fmt.Printf("Today we have %d sites to scan!\n", len(urls))
	storage := make(map[string]string)
	collectData(urls, depth, storage)

	// получаем значение поисковой строки из аргумента
	// и проставляем флаг, если строка была получена именно таким способом
	var recvCmdArg = false
	search := parseSearchPhrase()
	if search != "" {
		recvCmdArg = true
	}

	for {
		// если поисковая строка не была получена из аргумента
		// или была указана пустой - запрашиваем ее у пользователя
		if search == "" {
			fmt.Print("Enter search phrase and press ENTER to search ")
			fmt.Print("(or just press ENTER to exit): ")
			fmt.Scanln(&search)

			// если пользователь ничего не указал - выходим из приложения
			if search == "" {
				break
			}
		}

		found := findRelatedPages(storage, search)

		fmt.Println("Pages found:")
		if len(found) > 0 {
			for _, link := range found {
				fmt.Printf("[%s] %s\n", storage[link], link)
			}
		} else {
			fmt.Println("Nothing.")
		}

		// если поисковая строка была получена из аргумента,
		// то дальнейший интерактив не требуется - выходим из приложения
		if recvCmdArg {
			break
		}

		// сбрасываем значение поисковой строки
		// для выдачи запроса на ее ввод в следующей итерации
		search = ""
	}
}

// parseSearchPhrase возвращает поисковое значение, переданное при запуске приложения
func parseSearchPhrase() string {
	searchPtr := flag.String("search", "", "word/phrase to search in page title")
	flag.Parse()

	if *searchPtr != "" {
		fmt.Printf("Your search phrase is: %s\n", *searchPtr)
	}

	return *searchPtr
}

// collectData накапливает в словаре ссылки/заголовки страниц с загруженных сайтов
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

// findRelatedPages возвращает список ссылок на страницы,
// в заголовках которых содержится искомая фраза
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
