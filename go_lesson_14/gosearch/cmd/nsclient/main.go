package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
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
		var word string
		fmt.Scanln(&word)
		if word == "" {
			return
		}
		word += "\n"

		// делаем поисковый запрос
		_, err := conn.Write([]byte(word))
		if err != nil {
			return
		}

		// считываем длину ответа
		var lenPart []byte
		for {
			lenPart, _, err = r.ReadLine()
			if err != nil {
				return
			}
			if string(lenPart) != "" {
				break
			}
		}
		len, err := strconv.Atoi(string(lenPart))
		if err != nil {
			return
		}

		// зная длину, считываем текст ответа
		textPart := make([]byte, len)
		_, err = r.Read(textPart)
		if err != nil {
			return
		}

		fmt.Println(string(textPart))
	}
}
