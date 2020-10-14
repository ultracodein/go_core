// Package fibo implements some functionality inspired by Fibonacci numbers
package fibo

// NumberAtPos returns Fibonacci number at specified position and empty error message
// (or returns zero and detailed error message)
func NumberAtPos(pos uint) (errorMsg string, number uint) {
	switch {
	case pos >= 0 && pos <= 1:
		return "", pos
	case pos >= 2 && pos <= 20:
		var _, x = NumberAtPos(pos - 2)
		var _, y = NumberAtPos(pos - 1)
		return "", x + y
	case pos > 20:
		return "Getting numbers at pos > 20 is not supported!", 0
	}

	return "Unknown runtime error!", 0
}
