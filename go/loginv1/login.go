package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaultx/db"
	"vaultx/errorcheck"

	"golang.org/x/crypto/bcrypt"
)

type response struct {
	Authentication bool
}

type Creds struct {
	Email    string `json:"mailid"`
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
			conn.QueryRow(context.TODO(), "select password from users where mailid=$1", creds.Email).Scan(&password)

			if bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.Password)) == nil {
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
