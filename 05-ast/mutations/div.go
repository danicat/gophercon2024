package main

import "errors"

var ErrDivideByZero = errors.New("cannot divide by zero")

func divide(dividend, divisor int) (int, error) {
	if divisor == 0 {
		// I'm a comment
		return 0, ErrDivideByZero
	}

	return dividend / divisor, nil
}
