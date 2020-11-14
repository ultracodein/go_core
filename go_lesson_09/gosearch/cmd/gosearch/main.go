package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gosearch/pkg/cache"
	"gosearch/pkg/cache/btcache"
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/scheduler"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btstore"
)

// Сервер Интернет-поисковика GoSearch.
type gosearch struct {
	engine    *engine.Service
	scanner   crawler.Interface
	index     index.Interface
	storage   storage.Interface
	cache     cache.Interface
	scheduler *scheduler.Service

	initDone      bool
	initError     error
	cacheFiles    map[string]string
	schedulerFile string
	sites         []string
	depth         int
}

func main() {
	server := new()
	server.tryInitWithCache()

	// если есть устаревшие сайты - сканируем их и обновляем данные
	expired := server.scheduler.ExpiredSites()
	if len(expired) > 0 {
		go func() {
			server.updateDataByScan(expired)
		}()
	}

	// дожидаемся завершения инициализации поискового движка
	for {
		if !server.initDone {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if server.initError != nil {
		fmt.Println("Ошибка инициализации поискового движка! Функционал поиска недоступен.")
		return
	}

	server.searchLoop()
}

// tryInitWithCache загружает индекс и хранилище из файлов и запускает поисковый движок на них
func (gs *gosearch) tryInitWithCache() {
	cachedIndex, cachedStorage, err := gs.loadCache()
	if err != nil {
		log.Println("При загрузке данных из кэша произошла ошибка:", err)
		return
	}

	gs.index = cachedIndex
	gs.storage = cachedStorage
	gs.engine = engine.New(gs.index, gs.storage)
	gs.initDone = true
}

// updateDataByScan обновляет индекс и хранилище, инициализирует поисковый движок,
// сохраняет обновленные данные и состояние планировщика в файлы
func (gs *gosearch) updateDataByScan(sites []string) {
	fmt.Println("Идет сканирование сайтов. Пожалуйста, подождите...")
	scanned := gs.updateIndexByScan(sites)
	gs.scheduler.UpdateHistory(scanned)

	err := gs.saveCache()
	if err != nil {
		log.Println("При сохранении данных в кэш произошла ошибка:", err)
		gs.initError, gs.initDone = err, true
		return
	}

	err = gs.scheduler.SaveTo(gs.schedulerFile)
	if err != nil {
		log.Println("При сохранении состояния планировщика произошла ошибка:", err)
		gs.initError, gs.initDone = err, true
		return
	}

	if gs.engine == nil {
		gs.engine = engine.New(gs.index, gs.storage)
	}
	gs.initDone = true
}

// updateIndexByScan выполняет сканирование сайтов и индексирование данных
func (gs *gosearch) updateIndexByScan(sites []string) (scanned []string) {
	scanned = make([]string, 0)
	id := 0
	for _, site := range sites {
		log.Println("Сайт:", site)
		data, err := gs.scanner.Scan(site, gs.depth)
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
			continue
		}
		scanned = append(scanned, site)
	}
	return scanned
}

// new создаёт объект и службы сервера и возвращает указатель на него.
func new() *gosearch {
	gs := gosearch{}
	gs.scanner = spider.New()
	gs.index = hash.New()
	gs.storage = btstore.New()
	gs.cacheFiles = map[string]string{
		"index":   "index.bin",
		"storage": "storage.bin",
	}
	gs.cache = btcache.New(gs.cacheFiles)
	gs.schedulerFile = "scheduler.bin"
	gs.scheduler = gs.initScheduler()
	gs.depth = 2
	return &gs
}

func (gs *gosearch) initScheduler() *scheduler.Service {
	s, err := scheduler.LoadFrom(gs.schedulerFile)
	if err != nil {
		sites := []string{"https://www.ixbt.com/", "https://habr.com/"}
		expdays := 3
		log.Println("Планировщик не найден (инициализирован начальными значениями).")
		return scheduler.New(sites, expdays)
	}
	log.Println("Планировщик найден и загружен.")
	return s
}

// loadCache загружает индекс и хранилище из кэша (если он есть)
func (gs *gosearch) loadCache() (index.Interface, storage.Interface, error) {
	_, err := os.Stat(gs.cacheFiles["index"])
	if err != nil {
		return nil, nil, err
	}
	_, err = os.Stat(gs.cacheFiles["storage"])
	if err != nil {
		return nil, nil, err
	}
	ind, str, err := gs.cache.Load()
	if err != nil {
		return nil, nil, err
	}
	return *ind, *str, nil
}

// saveCache сохраняет индекс и хранилище в кэш
func (gs *gosearch) saveCache() error {
	return gs.cache.Create(&gs.index, &gs.storage)
}

// searchLoop выполняет поиск документов
func (gs *gosearch) searchLoop() {
	for {
		// запрашиваем у пользователя слово для поиска
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

		// сбрасываем значение для выдачи запроса на ввод слова в следующей итерации
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
