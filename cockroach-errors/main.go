package main

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
)

func main() {
	fmt.Println("Hello, World!")

	e1 := errors.New("error 1")
	e2 := errors.New("error 2")

	multierr.Combine(e1, e2)
}
