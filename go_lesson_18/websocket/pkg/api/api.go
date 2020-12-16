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
	messages []string
	pool     []*websocket.Conn
}

// New - конструктор API
func New(router *mux.Router, upgrader websocket.Upgrader, secret string) *API {
	api := API{
		router:   router,
		upgrader: upgrader,
		secret:   secret,
		messages: []string{},
		pool:     []*websocket.Conn{},
	}
	return &api
}

var mtx sync.Mutex

// Endpoints регистрирует конечные точки API
func (api *API) Endpoints() {
	api.router.HandleFunc("/send", api.sendHandler)
	api.router.HandleFunc("/messages", api.messagesHandler)
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
	conn.WriteMessage(websocket.TextMessage, []byte("OK"))

	// читаем сообщение
	_, msg, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	// добавляем полученное сообщение в кеш для отправки
	mtx.Lock()
	api.messages = append(api.messages, string(msg))
	mtx.Unlock()
}

func (api *API) messagesHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// добавляем соединение с клиентом в общий пул соединений
	mtx.Lock()
	api.pool = append(api.pool, conn)
	mtx.Unlock()

	// если клиент отключился - удаляем его соединение из пула
	defer func() {
		mtx.Lock()
		for i := range api.pool {
			if api.pool[i] == conn {
				api.pool = append(api.pool[:i], api.pool[i+1:]...)
				break
			}
		}
		mtx.Unlock()
	}()

	for {
		// если отправлять нечего - не блокируем кеш сообщений и пул соединений
		if len(api.messages) == 0 {
			continue
		}

		// иначе отправляем сообщения во все соединения
		mtx.Lock()
		for i := range api.pool {
			for j := range api.messages {
				api.pool[i].WriteMessage(websocket.TextMessage, []byte(api.messages[j]))
			}
		}
		// очищаем кеш сообщений
		api.messages = []string{}
		mtx.Unlock()
	}
}
