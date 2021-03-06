package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/rpcsrv"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/memstore"
)

// Сервер Интернет-поисковика GoSearch.
type gosearch struct {
	rpc     *rpcsrv.Service
	engine  *engine.Service
	scanner crawler.Interface
	index   index.Interface
	storage storage.Interface
	sites   []string
	depth   int
}

const maxThreads int = 10
const addr string = ":8000"

func main() {
	server := new()

	err := server.init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Поисковый движок инициализирован.")

	lsnr, err := server.start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("RPC-сервер запущен.")

	err = http.Serve(lsnr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// new создаёт объект и службы сервера и возвращает указатель на него.
func new() *gosearch {
	gs := gosearch{}
	gs.scanner = spider.New(maxThreads)
	gs.index = hash.New()
	gs.storage = memstore.New()
	gs.engine = engine.New(gs.index, gs.storage)
	gs.rpc = rpcsrv.New(gs.engine)
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

// saveData выполняет индексирование и сохранение документов в хранилище
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

// start регистрирует RPC-сервер и запускает его на указанном адресе
func (gs *gosearch) start() (net.Listener, error) {
	err := rpc.Register(gs.rpc)
	if err != nil {
		return nil, err
	}
	rpc.HandleHTTP()
	lsnr, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}
	return lsnr, nil
}
