package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	err := fmt.Errorf("simple error")

	err = errors.Wrap(err, "some error")
	fmt.Println(err)
}
