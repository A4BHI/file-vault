package main

import (
	"net/http"
	"vaultx/registerv1"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("../html")))
	http.HandleFunc("/register", registerv1.Register)
	http.ListenAndServe(":8080", nil)

}
