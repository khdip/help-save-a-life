package main

import (
	"log"

	"github.com/khdip/help-save-a-life/server/storage/postgres"
)

func main() {
	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
