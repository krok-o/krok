package main

import (
	"log"

	"github.com/krok-o/krok/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
