package main

import (
	"net/http"
	"vaultx/registerv1"
	"vaultx/verify"
)

type Test struct {
	H string
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../html")))
	http.HandleFunc("/register", registerv1.Register)
	http.HandleFunc("/verifyotp", verify.Verify)
	http.ListenAndServe(":8080", nil)

}
