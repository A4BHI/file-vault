package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaultx/db"
	"vaultx/email"
	"vaultx/errorcheck"
	"vaultx/masterkeys"
	"vaultx/session"

	"golang.org/x/crypto/bcrypt"
)

type response struct {
	Authentication bool `json:"auth"`
	Email          bool `json:"send_mail"`
}

type Creds struct {
	Mailid   string `json:"mailid"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("+++++++++++LOGIN++++++++++++++")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error from Login.go login function: ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		creds := Creds{}
		res := response{}
		json.Unmarshal(body, &creds)

		checklog := func() {
			conn, err := db.Connect()
			errorcheck.PrintError("Error connecting to db in login function", err)
			var password string
			var mailid string

			var username string
			conn.QueryRow(context.TODO(), "select hashed_passed,mailid,username from users where mailid=$1", creds.Mailid).Scan(&password, &mailid, &username)
			fmt.Println("check1")
			fmt.Println(password)
			fmt.Println(creds.Mailid)
			fmt.Println(creds.Password)
			fmt.Println(mailid)
			err = bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.Password))
			if err == nil {
				masterkeys.StorePassword(mailid, creds.Password)

				fmt.Println("check2")
				err := email.SendMail(mailid, username)
				if err != nil {
					fmt.Println("Error in Login()", err)
					res.Email = false
					json.NewEncoder(w).Encode(res)
					return
				}

				errorcheck.PrintError("error connecting to db in login.go ", err)

				session.SetSession(w, username, mailid, "nil", "login")
				fmt.Println(username)
				res.Email = true
				res.Authentication = true
				json.NewEncoder(w).Encode(res)

			} else {
				fmt.Println("check3")
				fmt.Println(err)
				res.Authentication = false
				json.NewEncoder(w).Encode(res)
				return
			}

		}

		checklog()

	}

}
