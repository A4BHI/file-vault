package verify

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaultx/errorcheck"
	"vaultx/otps"
	"vaultx/registerv1"
)

type userotp struct {
	Uotp string `json:"otp"`
}
type jsonresponse struct {
	Verified        bool `json:"Verified"`
	Account_Created bool `json:"Account_Created"`
}

func getSession(r *http.Request) string {
	cookie, err := r.Cookie("mail")
	if err != nil {
		fmt.Println("Error in otpverification.go error getting session from cookie", err)
	}

	b, err := hex.DecodeString(cookie.Value)

	if err != nil {
		fmt.Println("Error decoding session value in otpverification.go:", err)
	}

	return string(b)
}

func Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

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
		res := jsonresponse{}

		if response {
			res.Verified = true
			otps.DeleteOtp(emailid)
			cokkie := &http.Cookie{
				Name:     "mail",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			}
			http.SetCookie(w, cokkie)

			if registerv1.SaveToDB() {
				res.Account_Created = true
			} else {
				res.Account_Created = false
			}

		} else {
			res.Account_Created = false
			res.Verified = false

		}
		w.Header().Set("Content-Type", "Application/JSON")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println("Error Marshaling in otpverification.go :", err)
		}

	}
}
