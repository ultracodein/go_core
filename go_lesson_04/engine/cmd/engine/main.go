package main

import (
	"flag"
	"fmt"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/index"
)

func main() {
	urls := []string{"https://habr.com/", "https://www.cnews.ru/"}
	const depth int = 2
	var spider = spider.New()

	fmt.Printf("Today we have %d sites to scan!\n", len(urls))
	storage := index.Fill(spider, urls, depth)

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

		fmt.Print(storage.Find(search))

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
