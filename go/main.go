package main

import (
	"net/http"
	"time"
	auth "vaultx/loginv1"
	"vaultx/registerv1"
	"vaultx/session"
	"vaultx/verify"
)

func main() {

	go func() {
		for {
			session.CleanSessionFromDB()
			time.Sleep(time.Minute * 10)
		}

	}()
	http.Handle("/", http.FileServer(http.Dir("../html")))
	http.HandleFunc("/register", registerv1.Register)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/verifyotp", verify.Verify)
	http.HandleFunc("/resendotp", verify.ResendOtp)
	http.ListenAndServe(":8080", nil)

}
