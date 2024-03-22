package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on :10000")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path[1:]
	fmt.Fprintf(w, "Echo: %s", message)
}
