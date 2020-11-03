package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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
	server.loadCache()
	go func() {
		server.init()
		server.saveCache()
	}()
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
	gs.engine = engine.New(gs.index, gs.storage)
	gs.sites = []string{"https://www.ixbt.com/", "https://habr.com/"}
	gs.depth = 2
	return &gs
}

// loadCache загружает индекс и хранилище из кэша (если он есть)
func (gs *gosearch) loadCache() {
	_, err := os.Stat(gs.cacheFiles["index"])
	if err != nil {
		log.Println("Кэш не найден.")
		return
	}
	_, err = os.Stat(gs.cacheFiles["storage"])
	if err != nil {
		log.Println("Кэш не найден.")
		return
	}

	log.Println("Загрузка данных из кэша.")
	ind, str, err := gs.cache.Load()
	if err != nil {
		log.Println(err)
	}

	gs.index = *ind
	gs.storage = *str
	gs.engine = engine.New(gs.index, gs.storage)
}

// saveCache сохраняет индекс и хранилище в кэш
func (gs *gosearch) saveCache() {
	log.Println("Сохранение данных в кэш.")
	err := gs.cache.Create(&gs.index, &gs.storage)
	if err != nil {
		log.Println(err)
	}
}

// init производит сканирование сайтов и индексирование данных.
func (gs *gosearch) init() {
	log.Println("Сканирование сайтов.")
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
		log.Println("Индексирование документов.")
		gs.index.Add(data)
		log.Println("Сохранение документов.")
		err = gs.storage.StoreDocs(data)
		if err != nil {
			log.Println("ошибка при добавлении документов в хранилище:", err)
			continue
		}
	}
}

// run выполняет поиск документов
func (gs *gosearch) searchLoop() {
	// получаем значение поисковой строки из аргумента
	// и проставляем флаг, если строка была получена именно таким способом
	var fromCmd = false
	word := parseSearchWord()
	if word != "" {
		fromCmd = true
	}

	for {
		// если поисковая строка не была получена из аргумента
		// или была указана пустой - запрашиваем ее у пользователя
		if word == "" {
			fmt.Print("Enter word and press ENTER: ")
			fmt.Scanln(&word)
			// если пользователь ничего не указал - выходим из приложения
			if word == "" {
				break
			}
		}

		docs := gs.engine.Search(word)
		text := genOutput(docs)
		fmt.Print(text)

		// если поисковая строка была получена из аргумента,
		// то дальнейший интерактив не требуется - выходим из приложения
		if fromCmd {
			break
		}

		// сбрасываем значение поисковой строки
		// для выдачи запроса на ее ввод в следующей итерации
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

// parseSearchWord возвращает слово для поиска, переданное в виде флага при запуске приложения
func parseSearchWord() string {
	searchPtr := flag.String("search", "", "word to search in page title")
	flag.Parse()
	return *searchPtr
}
