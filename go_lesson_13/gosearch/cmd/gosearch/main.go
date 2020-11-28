package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

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

	cacheFiles    map[string]string
	schedulerFile string
	sites         []string
	depth         int
}

func main() {
	server := new()
	searchStarted := false
	ch := make(chan error)
	var wg sync.WaitGroup

	// пытаемся инициализировать и запустить сервер на кэше
	err := server.initWithCache()
	if err == nil {
		wg.Add(1)
		go func() {
			server.searchLoop(&wg)
		}()
		searchStarted = true
	} else {
		log.Println("При загрузке данных из кэша произошла ошибка:", err)
	}

	// если есть устаревшие сайты - сканируем их и обновляем данные сервера
	expired := server.scheduler.ExpiredSites()
	if len(expired) > 0 {
		go func() {
			server.updateDataByScan(expired, ch)
		}()
		err := <-ch
		if err != nil {
			log.Println("При обновлении данных поисковика произошла ошибка:", err)
			wg.Wait()
			return
		}

		if !searchStarted {
			wg.Add(1)
			server.searchLoop(&wg)
		}
	}

	wg.Wait()
}

// initWithCache загружает индекс и хранилище из файлов и запускает поисковый движок на них
func (gs *gosearch) initWithCache() error {
	cachedIndex, cachedStorage, err := gs.loadCache()
	if err != nil {
		return err
	}

	gs.index = cachedIndex
	gs.storage = cachedStorage
	gs.engine = engine.New(gs.index, gs.storage)
	log.Println("Поисковый движок запущен на кэше.")
	return nil
}

// updateDataByScan обновляет индекс и хранилище, инициализирует поисковый движок,
// сохраняет обновленные данные и состояние планировщика в файлы
func (gs *gosearch) updateDataByScan(sites []string, ch chan error) {
	fmt.Println("Идет сканирование сайтов. Пожалуйста, подождите...")
	scanned := gs.updateIndexByScan(sites)

	err := gs.saveCacheAndHistory(scanned)
	if err != nil {
		ch <- err
		return
	}

	if gs.engine == nil {
		gs.engine = engine.New(gs.index, gs.storage)
	}
	ch <- nil
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

func (gs *gosearch) saveCacheAndHistory(scanned []string) error {
	err := gs.saveCache()
	if err != nil {
		return err
	}

	gs.scheduler.UpdateHistory(scanned)
	err = gs.scheduler.SaveTo(gs.schedulerFile)
	if err != nil {
		return err
	}

	return nil
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
		sites := []string{"https://www.ixbt.com/", "https://cnews.ru/"}
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
func (gs *gosearch) searchLoop(wg *sync.WaitGroup) {
	defer wg.Done()
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
