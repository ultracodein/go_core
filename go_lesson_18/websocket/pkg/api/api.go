package api

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// API реализует программный интерфейс сервера обмена сообщениями
type API struct {
	router   *mux.Router
	upgrader websocket.Upgrader
	secret   string
	msgs     chan string
	pool     []chan string
}

// New - конструктор API
func New(router *mux.Router, upgrader websocket.Upgrader, secret string) *API {
	api := API{
		router:   router,
		upgrader: upgrader,
		secret:   secret,
		msgs:     make(chan string),
		pool:     []chan string{},
	}
	return &api
}

var mtx sync.Mutex

// Endpoints регистрирует конечные точки API
func (api *API) Endpoints() {
	api.router.HandleFunc("/send", api.sendHandler)
	api.router.HandleFunc("/messages", api.messagesHandler)
}

// ForwardMessages направляет сообщения из общего канала в каналы пула
func (api *API) ForwardMessages() {
	for msg := range api.msgs {
		for _, ch := range api.pool {
			ch <- msg
		}
	}
}

func (api *API) sendHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// читаем пароль и проверяем его
	_, secret, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	if string(secret) != api.secret {
		http.Error(w, "указан неверный пароль", http.StatusUnauthorized)
		return
	}

	// отправляем подтверждение
	err = conn.WriteMessage(websocket.TextMessage, []byte("OK"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// читаем сообщение
	_, msg, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	// пишем полученное сообщение в общий канал
	api.msgs <- string(msg)
}

func (api *API) messagesHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// добавляем канал для отправки сообщений клиенту в общий пул
	mtx.Lock()
	ch := make(chan string)
	api.pool = append(api.pool, ch)
	mtx.Unlock()

	defer deleteChan(&api.pool, ch)

	// отправка сообщений клиенту (при их наличии в канале)
	for msg := range ch {
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func deleteChan(chPtr *[]chan string, ch chan string) {
	mtx.Lock()
	pool := *chPtr
	for i := range pool {
		if pool[i] == ch {
			pool = append(pool[:i], pool[i+1:]...)
			break
		}
	}
	mtx.Unlock()
}
