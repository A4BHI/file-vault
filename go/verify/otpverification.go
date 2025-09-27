package verify

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaultx/errorcheck"
	"vaultx/otps"
)

type userotp struct {
	Uotp string `json:"otp"`
}

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
		o := userotp{}
		err = json.Unmarshal(UserInputOtp, &o)
		errorcheck.Nigger("Error unmarshalling json File:otpverification.go function:Verify", err)
		emailid := getSession(r)
		fmt.Println("Emailid:", emailid)
		fmt.Println("User Input otp:", o.Uotp)
		a := otps.GetOtp(emailid)
		fmt.Println("Generated Otp", a)
		response := otps.CompareOtp(emailid, o.Uotp)

		if response {
			otps.DeleteOtp(emailid)
			cokkie.MaxAge = -1
			fmt.Fprintf(w, `{"verified":true}`)
		} else {
			fmt.Fprintf(w, `{"verified":false}`)
		}

	}
}
