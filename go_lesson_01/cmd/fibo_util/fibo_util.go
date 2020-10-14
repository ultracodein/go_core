package main

import (
	"fmt"
	"go_core/go_lesson_01/fibo"
)

func main() {
	positions := []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 20, 21}

	for _, pos := range positions {
		var errorMsg, number = fibo.NumberAtPos(pos)

		if errorMsg == "" {
			fmt.Printf("F[%d] = %d\n", pos, number)
		} else {
			fmt.Printf("F[%d] = Oops! Something went wrong. Details: %s\n", pos, errorMsg)
		}
	}
}
