package main

import (
	"github.com/joansais/go-tutorials/wiki"
	"log"
	"flag"
)

const (
	DEFAULT_ASSETS_DIR = "assets"
	DEFAULT_STORAGE_DIR = "data/pages"
)

func main() {
	assetsDir := flag.String("assets", DEFAULT_ASSETS_DIR, "location of HTML templates")
	storageDir := flag.String("storage", DEFAULT_STORAGE_DIR, "storage directory for wiki pages")
	flag.Parse()

	store := wiki.NewDiskStore(*storageDir)
	syntax := wiki.NewMarkdownSyntax(store)
	server := wiki.NewServer(store, syntax, *assetsDir)
	err := server.Start(":8080")
	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
}
