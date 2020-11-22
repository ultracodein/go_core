package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type action struct {
	actor, command string
}

func main() {
	N := 11
	results := make(map[string]int)

	// играем N раундов
	for i := 1; i <= N; i++ {
		winner := playRound()
		fmt.Printf("--- Winner: %s---\n\n", winner)
		results[winner]++
	}

	// выводим итоговый счет
	for player, score := range results {
		fmt.Printf("%s: %d\n", player, score)
	}
}

func playRound() string {
	// инициализируем игроков
	var wg sync.WaitGroup
	ch := make(chan action)
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go playerLoop(i, ch, &wg)
	}

	// даем команду к началу раунда
	ch <- action{
		actor:   "Referee",
		command: "begin",
	}

	winner := ""
	for {
		msg := <-ch

		// при появлении команды stop фиксируем победителя раунда
		if msg.command == "stop" {
			winner = msg.actor
			close(ch)
			break
		} else {
			// на другие команды не реагируем
			ch <- msg
		}
	}
	wg.Wait()

	return winner
}

func playerLoop(id int, ch chan action, wg *sync.WaitGroup) {
	defer wg.Done()

	// инициализируем игрока (по умолчанию он - "принимающий")
	name := "Player " + strconv.Itoa(id)
	pinger := false
	reply := "pong"

	for {
		msg, open := <-ch
		if !open {
			return
		}

		// выполняем подачу и запоминаем, что мы - "подающие"
		if msg.command == "begin" {
			fmt.Printf("--- Serving: %s ---\n", name)

			pinger = true
			reply = "ping"
			if goodKick(name, reply, ch) {
				return
			}
			continue
		}

		// отвечаем на удары соперника
		if (msg.command == "ping" && !pinger) || (msg.command == "pong" && pinger) {
			if goodKick(name, reply, ch) {
				return
			}
			continue
		}

		// на собственные удары и stop не реагируем
		ch <- msg
	}
}

func goodKick(actor, command string, ch chan action) bool {
	if goodKickTime() {
		doAction(actor, "stop", ch)
		return true
	}

	doAction(actor, command, ch)
	return false
}

func goodKickTime() bool {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5)
	if n != 0 {
		return false
	}

	return true
}

func doAction(actor, command string, ch chan action) {
	fmt.Printf("%s: %s\n", actor, command)

	ch <- action{
		actor:   actor,
		command: command,
	}
}
