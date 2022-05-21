package main

import (
	"fmt"
	"gii"
	"log"
	"net/http"
)

func main() {
	engine := gii.New()

	engine.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})

	engine.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!\n")
	})

	log.Fatal(engine.Run(":8080"))
}
