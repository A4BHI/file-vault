package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaultx/db"
	"vaultx/errorcheck"
	auth "vaultx/loginv1"
	"vaultx/masterkeys"
	"vaultx/otps"
	"vaultx/registerv1"
	"vaultx/session"
)

type userotp struct {
	Uotp string `json:"otp"`
}

//	type loginresponse struct {
//		Login bool `json:"login"`
//	}
type jsonresponse struct {
	Login           bool   `json:"login"`
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

func GetUserName(mailid string) string {
	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting db in GetUserName() in otpverification.go", err)
	var username string
	conn.QueryRow(context.TODO(), "Select username from users where mailid=$1", mailid).Scan(&username)
	return username
}

func Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// log := loginresponse{}
		res := jsonresponse{}
		UserInputOtp, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error getting otp from frontend:", err)
		}
		o := userotp{}
		err = json.Unmarshal(UserInputOtp, &o)
		errorcheck.PrintError("Error unmarshalling json File:otpverification.go function:Verify", err)

		sessionid := getSession(r)
		_, mailid, _, _, session_type, ok := session.GetSession(sessionid)
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

			if session_type == "login" {
				// log.Login = true
				password, ok := masterkeys.GetPassword(mailid)
				if !ok {
					fmt.Println("Error getting password from getPassword() in otpverification")
					return
				}
				m := masterkeys.GenerateKey(mailid, password.(string))
				fmt.Println("MasterKey:", m)
				masterkeys.DeletePassword(mailid)
				auth.Setjwtkey(w, GetUserName(mailid)) //jwt ivide
				res.Login = true
				res.Verified = true
				json.NewEncoder(w).Encode(res)
				if session.DeleteSession(sessionid) {
					fmt.Println("PASS")
				} else {
					fmt.Println("FAIL")
				}
				return
			} else if session_type == "register" {
				if registerv1.SaveToDB(sessionid) {
					res.Account_Created = true
					password, ok := masterkeys.GetPassword(mailid)
					if !ok {
						fmt.Println("Error getting password from getPassword() in otpverification")
						return
					}
					m := masterkeys.GenerateKey(mailid, password.(string))
					fmt.Println("MasterKey:", m)
					masterkeys.DeletePassword(mailid)
				} else {
					masterkeys.DeletePassword(mailid)
					res.Account_Created = false
					conn, _ := db.Connect()
					conn.Exec(context.TODO(), "Delete from sessions where sessionid=$1", sessionid)
				}
			}

		} else {
			res.Account_Created = false
			res.Verified = false
			masterkeys.DeletePassword(mailid)

		}
		w.Header().Set("Content-Type", "Application/JSON")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println("Error Marshaling in otpverification.go :", err)
		}

	}
}
