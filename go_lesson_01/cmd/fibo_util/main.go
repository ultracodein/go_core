package main

import (
	"flag"
	"fmt"
	"go_core/go_lesson_01/fibo"
)

func main() {
	modePtr := flag.String("mode", "unknown", "launch mode: calc or demo")
	posPtr := flag.Uint("pos", 0, "pos of Fibonacci number")
	flag.Parse()

	switch {
	case *modePtr == "calc":
		printFiboNums(*posPtr)
	case *modePtr == "demo":
		showDemo()
	default:
		fmt.Println("Please specify launch mode: -mode=calc or -mode=demo")
	}
}

func printFiboNums(pos uint) {
	var num, err = fibo.Num(pos)

	if err == nil {
		fmt.Printf("F[%d] = %d\n", pos, num)
	} else {
		fmt.Printf("F[%d] = Oops! Something went wrong. Details: %s\n", pos, err.Error())
	}
}

func showDemo() {
	positions := []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 20, 21}

	for _, pos := range positions {
		printFiboNums(pos)
	}
}
