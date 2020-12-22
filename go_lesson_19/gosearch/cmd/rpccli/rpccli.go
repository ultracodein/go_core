package main

import (
	"bufio"
	"fmt"
	"gosearch/pkg/crawler"
	"log"
	"net/rpc"
	"os"
	"strings"
)

const srvaddr string = "localhost:8000"

func main() {
	// подключаемся к серверу
	client, err := rpc.DialHTTP("tcp4", srvaddr)
	if err != nil {
		log.Fatal(err)
	}

	// поиск в интерактивном режиме
	search(client)
}

func search(client *rpc.Client) {
	r := bufio.NewReader(os.Stdin)
	for {
		// получаем поисковую фразу
		word, err := r.ReadString('\n')
		if err != nil {
			log.Fatalf("не удалось прочитать поисковую фразу из stdin: %v", err)
		}
		word = strings.TrimSuffix(word, "\n")
		if word == "" {
			return
		}

		// выполняем запрос по RPC
		var docs []crawler.Document
		err = client.Call("Service.Search", word, &docs)
		if err != nil {
			log.Fatal("ошибка при вызове метода RPC:", err)
		}

		// обрабатываем результат
		if len(docs) == 0 {
			fmt.Println("По вашему запросу ничего не найдено.")
			continue
		}
		res := searchResult(docs)
		fmt.Print(res)
	}
}

func searchResult(docs []crawler.Document) string {
	var result string
	for _, doc := range docs {
		result += fmt.Sprintf("[%d] %s\n", doc.ID, doc.Title)
	}
	return result
}
