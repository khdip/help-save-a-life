package main

import (
	"help-save-a-life/server/storage/postgres"
	"log"
)

func main() {
	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
