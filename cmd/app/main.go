package main

import (
	"log"

	"github.com/iktkhor/url-downloader/internal/pkg/app"
)

func main() {
	a := app.New()

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
