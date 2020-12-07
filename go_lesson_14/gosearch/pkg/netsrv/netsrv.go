package netsrv

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"

	"gosearch/pkg/crawler"
	"gosearch/pkg/engine"
)

// Service - сетевая служба
type Service struct {
	addr   string
	engine *engine.Service
}

// New - конструктор.
func New(addr string, engine *engine.Service) *Service {
	s := Service{
		addr:   addr,
		engine: engine,
	}
	return &s
}

// Start запускает сетевую службу
func (s *Service) Start() error {
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handler(conn)
	}

	return nil
}

// handler обрабатывает запросы пользователей (выполняет поиск сайтов)
func (s *Service) handler(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		conn.SetDeadline(time.Now().Add(time.Second * 10))
		msg, _, err := r.ReadLine()
		if err != nil {
			return
		}

		conn.SetDeadline(time.Now().Add(time.Second * 10))
		word := string(msg)
		docs := s.engine.Search(word)
		data := reply(docs)
		fmt.Fprintf(conn, data+"<")
		conn.SetDeadline(time.Now().Add(time.Second * 10))
	}
}

// reply формирует ответ (результат поиска)
func reply(docs []crawler.Document) string {
	if docs == nil {
		return "Nothing."
	}

	var text string = ""
	for _, doc := range docs {
		text += "[" + strconv.Itoa(doc.ID) + "] " + doc.Title + "; "
	}
	return text
}
