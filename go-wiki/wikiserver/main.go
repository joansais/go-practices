package main

import (
	"github.com/joansais/go-tutorials/go-wiki/wiki"
	"net/http"
	"log"
)

func main() {
	wiki.RegisterServices()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
}
