package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp4", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Подключение установлено. Введите поисковую фразу:")
	searchLoop(conn)
}

func searchLoop(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		// получаем слово для поиска
		word := ""
		fmt.Scanln(&word)
		if word == "" {
			return
		}

		// делаем поисковый запрос
		fmt.Fprintf(conn, word+"\n")

		// считываем ответ
		reply, err := r.ReadString('<')
		if err != nil {
			return
		}

		fmt.Println(string(reply))
	}
}
