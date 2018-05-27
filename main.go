// Golang-HTTP-WebServer project main.go
// See https://golang.org/doc/articles/wiki/#tmp_3

package main

import (
	"fmt"
	"log"
	"net/http"
)

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You asked for %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", HttpHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
