package main

import (
	"flag"
	"fmt"
	"go_core/go_lesson_01/fibo"
)

func printFiboNumber(pos uint) {
	var errorMsg, number = fibo.NumberAtPos(pos)

	if errorMsg == "" {
		fmt.Printf("F[%d] = %d\n", pos, number)
	} else {
		fmt.Printf("F[%d] = Oops! Something went wrong. Details: %s\n", pos, errorMsg)
	}
}

func showDemo() {
	positions := []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 20, 21}

	for _, pos := range positions {
		printFiboNumber(pos)
	}
}

func main() {
	modePtr := flag.String("mode", "unknown", "launch mode: calc or demo")
	posPtr := flag.Uint("pos", 0, "pos of Fibonacci number")
	flag.Parse()

	switch {
	case *modePtr == "calc":
		printFiboNumber(*posPtr)
	case *modePtr == "demo":
		showDemo()
	default:
		fmt.Println("Please specify launch mode: -mode=calc or -mode=demo")
	}
}
