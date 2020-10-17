// Package fibo implements some functionality inspired by Fibonacci numbers
package fibo

import "errors"

// Num returns Fibonacci number at specified position and nil error
// (or returns zero and error with detailed message)
func Num(pos uint) (num uint, err error) {
	switch {
	case pos >= 0 && pos <= 1:
		return pos, nil
	case pos >= 2 && pos <= 20:
		var x, _ = Num(pos - 2)
		var y, _ = Num(pos - 1)
		return x + y, nil
	case pos > 20:
		return 0, errors.New("getting numbers at pos > 20 is not supported")
	}

	return 0, errors.New("unknown runtime error")
}
