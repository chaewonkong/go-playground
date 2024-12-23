package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	m := make(map[string]string)
	m["id"] = "cwkong"
	m["pw"] = "1234"

	keys := maps.Keys(m)

	for k := range keys {
		fmt.Println(k)
	}

	arr := slices.Collect(keys)
	fmt.Println("print")
	_ = arr
}
