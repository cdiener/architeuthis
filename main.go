package main

import (
	"fmt"
	"log"
)

func main() {
	m, err := Summarize("test.k2")
	if err != nil {
		log.Fatal(err)
	}
	k := "3864"
	v := m[k]
	fmt.Printf("%s: %v\n", k, *v)
}
