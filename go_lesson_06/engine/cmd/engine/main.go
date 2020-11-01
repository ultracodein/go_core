package main

import (
	"flag"
	"fmt"
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/index"
	"strconv"
)

func main() {
	urls := []string{"https://habr.com/", "https://www.cnews.ru/"}
	const depth int = 2

	// загружаем данные в индекс
	var index = index.New()
	var spider = spider.New()
	err := index.Fill(spider, urls, depth)
	if err != nil {
		return
	}

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
			fmt.Print("Enter word and press ENTER: ")
			fmt.Scanln(&search)

			// если пользователь ничего не указал - выходим из приложения
			if search == "" {
				break
			}
		}

		docs := index.Find(search)
		text := genResult(docs)
		fmt.Print(text)

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

func genResult(docs []*crawler.Document) string {
	if docs == nil {
		return "Nothing.\n"
	}

	var text string = ""
	for _, doc := range docs {
		text += "[" + strconv.Itoa(doc.ID) + "] " + doc.Title + "\n"
	}
	return text
}

// parseSearchPhrase возвращает поисковое значение, переданное при запуске приложения
func parseSearchPhrase() string {
	searchPtr := flag.String("search", "", "word to search in page title")
	flag.Parse()
	return *searchPtr
}
