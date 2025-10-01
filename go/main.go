package main

import (
	"net/http"
	"vaultx/registerv1"
	"vaultx/verify"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("../html")))
	http.HandleFunc("/register", registerv1.Register)
	http.HandleFunc("/verifyotp", verify.Verify)
	http.HandleFunc("/resendotp", verify.ResendOtp)
	http.ListenAndServe(":8080", nil)

}
