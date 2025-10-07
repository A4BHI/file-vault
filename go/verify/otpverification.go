package verify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaultx/errorcheck"
	"vaultx/otps"
	"vaultx/registerv1"
	"vaultx/session"
)

type userotp struct {
	Uotp string `json:"otp"`
}
type jsonresponse struct {
	Verified        bool   `json:"Verified"`
	Account_Created bool   `json:"Account_Created"`
	Internal_Error  string `json:"Internal_Error"`
}

func getSession(r *http.Request) string {
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		fmt.Println("Error in otpverification.go error getting session from cookie", err)
	}

	b := cookie.Value

	if err != nil {
		fmt.Println("Error decoding session value in otpverification.go:", err)
	}

	return b
}

func Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		res := jsonresponse{}
		UserInputOtp, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error getting otp from frontend:", err)
		}
		o := userotp{}
		err = json.Unmarshal(UserInputOtp, &o)
		errorcheck.PrintError("Error unmarshalling json File:otpverification.go function:Verify", err)
		sessionid := getSession(r)
		_, mailid, _, _, ok := session.GetSession(sessionid)
		if !ok {
			res.Verified = false
			res.Account_Created = false
			res.Internal_Error = "Cannot retrieve mailid from the sessionid"
			json.NewEncoder(w).Encode(res)
			return

		}

		fmt.Println("Emailid:", mailid)
		fmt.Println("User Input otp:", o.Uotp)
		a := otps.GetOtp(mailid)
		fmt.Println("Generated Otp", a)
		response := otps.CompareOtp(mailid, o.Uotp)

		if response {
			res.Verified = true
			otps.DeleteOtp(mailid)
			cokkie := &http.Cookie{
				Name:     "sessionid",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			}
			http.SetCookie(w, cokkie)

			if registerv1.SaveToDB(sessionid) {
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
