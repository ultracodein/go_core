package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gosearch/pkg/cache"
	"gosearch/pkg/cache/bstcache"
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/bststore"
)

// Сервер Интернет-поисковика GoSearch.
type gosearch struct {
	engine  *engine.Service
	scanner crawler.Interface
	index   index.Interface
	storage storage.Interface
	cache   cache.Interface

	cacheFiles map[string]string
	sites      []string
	depth      int
}

func main() {
	server := new()

	err := server.loadCache()
	if err != nil {
		log.Printf("Данные из кэша не загружены: [%s]", err)
	} else {
		log.Println("Данные загружены из кэша.")
	}

	go func() {
		log.Println("Сканирование сайтов.")
		server.scan()
		log.Println("Сохранение данных в кэш.")
		server.saveCache()
	}()

	// даем запуститься рутине, чтобы ее вывод не смешался с выводом цикла поиска
	time.Sleep(2 * time.Second)
	server.searchLoop()
}

// new создаёт объект и службы сервера и возвращает указатель на него.
func new() *gosearch {
	gs := gosearch{}
	gs.scanner = spider.New()
	gs.index = hash.New()
	gs.storage = bststore.New()
	gs.cacheFiles = map[string]string{
		"index":   "index.bin",
		"storage": "storage.bin",
	}
	gs.cache = bstcache.New(gs.cacheFiles)
	gs.sites = []string{"https://www.ixbt.com/", "https://habr.com/"}
	gs.depth = 2
	return &gs
}

// loadCache загружает индекс и хранилище из кэша (если он есть)
func (gs *gosearch) loadCache() error {
	_, err := os.Stat(gs.cacheFiles["index"])
	if err != nil {
		return err
	}
	_, err = os.Stat(gs.cacheFiles["storage"])
	if err != nil {
		return err
	}
	ind, str, err := gs.cache.Load()
	if err != nil {
		return err
	}
	gs.engine = engine.New(*ind, *str)
	return nil
}

// saveCache сохраняет индекс и хранилище в кэш
func (gs *gosearch) saveCache() {
	err := gs.cache.Create(&gs.index, &gs.storage)
	if err != nil {
		log.Println(err)
	}
}

// scan производит сканирование сайтов и индексирование данных.
func (gs *gosearch) scan() {
	id := 0
	for _, url := range gs.sites {
		log.Println("Сайт:", url)
		data, err := gs.scanner.Scan(url, gs.depth)
		if err != nil {
			continue
		}
		for i := range data {
			data[i].ID = id
			id++
		}
		gs.index.Add(data)
		err = gs.storage.StoreDocs(data)
		if err != nil {
			log.Printf("Ошибка при добавлении документов в хранилище: [%s]", err)
			continue
		}
	}
	gs.engine = engine.New(gs.index, gs.storage)
}

// run выполняет поиск документов
func (gs *gosearch) searchLoop() {
	for {
		if gs.engine == nil {
			time.Sleep(10 * time.Second)
			continue
		}

		var word string
		fmt.Print("Enter word and press ENTER: ")
		fmt.Scanln(&word)
		// если пользователь ничего не указал - выходим из приложения
		if word == "" {
			break
		}

		docs := gs.engine.Search(word)
		text := genOutput(docs)
		fmt.Print(text)

		// сбрасываем значение поисковой строки для выдачи запроса на ее ввод в следующей итерации
		word = ""
	}
}

func genOutput(docs []crawler.Document) string {
	if docs == nil {
		return "Nothing.\n"
	}

	var text string = ""
	for _, doc := range docs {
		text += "[" + strconv.Itoa(doc.ID) + "] " + doc.Title + "\n"
	}
	return text
}
