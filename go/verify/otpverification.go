package verify

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"vaultx/otps"
)

func getSession(r *http.Request) string {
	cokkie, err := r.Cookie("mail")
	if err != nil {
		fmt.Println("Error in otpverification.go error getting session from cookie", err)
	}

	b, err := hex.DecodeString(cokkie.Value)

	if err != nil {
		fmt.Println("Error decoding session value in otpverification.go:", err)
	}

	return string(b)
}

func Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cokkie, err := r.Cookie("mail")
		if err != nil {
			fmt.Println("Error fetching cookie File:Otpverification.go function Verify():", err)
		}
		UserInputOtp, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error getting otp from frontend:", err)
		}

		emailid := getSession(r)

		// otp := otps.GetOtp(emailid)

		response := otps.CompareOtp(emailid, string(UserInputOtp))

		if response {
			otps.DeleteOtp(emailid)
			cokkie.MaxAge = -1
		} else {

		}

	}
}
