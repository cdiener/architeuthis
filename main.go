package main

import (
	"log"

	"github.com/cdiener/architeuthis/cmd"
)

func main() {
	m, err := cmd.Summarize("positive_1.k2")
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.SaveMapping(m, "test.csv")
	if err != nil {
		log.Fatal(err)
	}
}
