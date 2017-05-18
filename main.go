package main

import (
	"log"

	"github.com/damoye/godo/downloader"
)

func main() {
	urlraw := "https://github.com/qt/qt/archive/4.8.zip"
	d, err := downloader.New(urlraw)
	if err != nil {
		log.Fatal(err)
	}
	err = d.Run()
	if err != nil {
		log.Fatal(err)
	}
}
