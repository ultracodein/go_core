package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btstore"
)

// Сервер Интернет-поисковика GoSearch.
type gosearch struct {
	engine  *engine.Service
	scanner crawler.Interface
	index   index.Interface
	storage storage.Interface
	sites   []string
	depth   int
}

const maxThreads int = 10

func main() {
	server := new()

	err := server.init()
	if err != nil {
		log.Fatal(err)
	}

	server.engine = engine.New(server.index, server.storage)
	fmt.Println("Поисковый движок запущен.")

	err = server.start()
	if err != nil {
		log.Fatal(err)
	}
}

// start запускает сетевую службу
func (gs *gosearch) start() error {
	listener, err := net.Listen("tcp4", ":8000")
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go gs.handler(conn)
	}

	return nil
}

// handler обрабатывает запросы пользователей (выполняет поиск сайтов)
func (gs *gosearch) handler(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		conn.SetDeadline(time.Now().Add(time.Second * 10))
		msg, _, err := r.ReadLine()
		if err != nil {
			return
		}

		conn.SetDeadline(time.Now().Add(time.Second * 5))
		word := string(msg)
		docs := gs.engine.Search(word)
		data := reply(docs)
		reply := fmt.Sprintf("%d\n%s\n", len(data), data)

		_, err = conn.Write([]byte(reply))
		if err != nil {
			return
		}
	}
}

// new создаёт объект и службы сервера и возвращает указатель на него.
func new() *gosearch {
	gs := gosearch{}
	gs.scanner = spider.New(maxThreads)
	gs.index = hash.New()
	gs.storage = btstore.New()
	gs.sites = []string{
		"https://www.hpe.com", "https://www.microsoft.com", "https://www.oracle.com", "https://www.citrix.com", "https://www.python.org",
		"http://cisco.com", "https://telegram.org", "https://www.tesla.com", "https://www.spacex.com", "https://www.formula1.com",
	}
	gs.depth = 2
	return &gs
}

// init инициирует сканирование и обработку полученных данных
func (gs *gosearch) init() error {
	docs, errs := gs.scanner.Scan(gs.sites, gs.depth)
	if len(errs) == len(gs.sites) {
		return errors.New("ни один из сайтов не отсканирован")
	}

	err := gs.saveData(docs)
	if err != nil {
		return err
	}

	return nil
}

// store выполняет индексирование и сохранение документов в хранилище
func (gs *gosearch) saveData(docs []crawler.Document) error {
	id := 0
	for i := range docs {
		docs[i].ID = id
		id++
	}

	gs.index.Add(docs)

	err := gs.storage.StoreDocs(docs)
	if err != nil {
		return err
	}

	return nil
}

// reply формирует ответ (результат поиска)
func reply(docs []crawler.Document) string {
	if docs == nil {
		return "Nothing.\n"
	}

	var text string = ""
	for _, doc := range docs {
		text += "[" + strconv.Itoa(doc.ID) + "] " + doc.Title + "\n"
	}
	return text
}
