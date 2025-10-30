package main

import (
	"net/http"
	"time"
	"vaultx/encryption"
	auth "vaultx/loginv1"
	"vaultx/registerv1"
	"vaultx/session"
	"vaultx/verify"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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
	http.HandleFunc("/checksession", session.Check)
	http.HandleFunc("/encrypt-file", encryption.Backend_Encryption)
	http.HandleFunc("/frontend-encrypt", encryption.Frontend_Enc)
	// http.HandleFunc("/jwt", auth.Setjwtkey)
	http.ListenAndServe(":8080", nil)

}
