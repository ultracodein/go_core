package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var address string = "ws://localhost:8000"
var secret string = "password"

func main() {
	go messages()
	send()
}

func send() {
	r := bufio.NewReader(os.Stdin)
	for {
		msg, _ := r.ReadString('\n')
		if msg == "\n" {
			return
		}

		// подключаемся
		conn, _, err := websocket.DefaultDialer.Dial(address+"/send", nil)
		if err != nil {
			log.Fatalf("не удалось подключиться к серверу: %v", err)
		}
		defer conn.Close()

		// передаем пароль
		err = conn.WriteMessage(websocket.TextMessage, []byte(secret))
		if err != nil {
			log.Fatalf("не удалось отправить пароль серверу: %v", err)
		}

		// проверяем готовность сервера к приему сообщений
		_, banner, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("не удалось прочитать ответ сервера: %v", err)
		}
		if string(banner) != "OK" {
			log.Fatal("ошибка авторизации")
		}

		// отправляем сообщение и отключаемся
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Fatalf("не удалось отправить сообщение серверу: %v", err)
		}
		conn.Close()
	}
}

func messages() {
	conn, _, err := websocket.DefaultDialer.Dial(address+"/messages", nil)
	if err != nil {
		log.Fatalf("не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("не удалось прочитать сообщение от сервера: %v", err)
		}
		if string(message) != "" {
			fmt.Printf("%s", message)
		}

		time.Sleep(time.Second)
	}
}
