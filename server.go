//go:build !wasm

package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed assets/*
var assets embed.FS

func main() {
	rootFS, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(rootFS)))
	log.Println("listening on port 7777")
	err = http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
