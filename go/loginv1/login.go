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

	"golang.org/x/crypto/bcrypt"
)

type response struct {
	Authentication bool
	Email          bool
}

type Creds struct {
	Username string `json:"mailid"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
			conn.QueryRow(context.TODO(), "select password,mailid from users where username=$1", creds.Username).Scan(&password, &mailid)

			if bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.Password)) == nil {
				err := email.SendMail(mailid, creds.Username)
				if err != nil {
					fmt.Println("Error in Login()", err)
					res.Email = false
					json.NewEncoder(w).Encode(res)
					return
				}

				res.Email = true
				res.Authentication = true

			} else {
				res.Authentication = false
				json.NewEncoder(w).Encode(res)
				return
			}

		}

		checklog()

	}

}
