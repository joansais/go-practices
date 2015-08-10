package main

import (
	"github.com/joansais/go-tutorials/go-wiki/wiki"
	"net/http"
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

	pageStore := wiki.NewDiskStore(*storageDir)
	wiki.SetAssetsDir(*assetsDir);
	wiki.SetPageStore(pageStore);
	wiki.SetSyntaxHandler(wiki.NewMarkdownSyntax(pageStore));
	wiki.RegisterServices()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
}
