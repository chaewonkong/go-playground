package main

import (
	"fmt"
	"sync"
)

func main() {
	m := sync.Map{}
	m.Store("key", 1)
	if v, ok := m.Load("key"); ok {
		fmt.Println(v)
	}
}
