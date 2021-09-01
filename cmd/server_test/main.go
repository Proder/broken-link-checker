package main

import (
	"broken-link-checker/app"
	"log"
)

func main() {
	if err := app.RunServerTest(); err != nil {
		log.Fatal(err)
	}
}
