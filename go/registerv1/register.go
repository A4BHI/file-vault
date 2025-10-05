package registerv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"vaultx/db"
	"vaultx/email"
	"vaultx/errorcheck"
	"vaultx/session"
)

type Creds struct {
	Username string `json:"username"`
	Email    string `json:"mailid"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("++++++++++Register++++++++++")
		body, err := io.ReadAll(r.Body)
		errorcheck.Nigger("Error From register.go in register Function:", err)

		var exists bool
		UserInfo := Creds{}
		err = json.Unmarshal(body, &UserInfo)
		errorcheck.Nigger("Error in register.go Register function json unmarshal", err)

		fmt.Println(UserInfo.Password)

		conn, err := db.Connect()
		errorcheck.Nigger("Error creating connection to psql (Register.go file in Register function)", err)

		row := conn.QueryRow(context.Background(), "select exists( select 1 from users where username = $1 or mailid=$2)", UserInfo.Username, UserInfo.Email)
		row.Scan(&exists)

		if exists {
			w.Header().Set("Content-Type", "Application/JSON")
			fmt.Fprintf(w, `{"ok":false}`) // send response as acc already registered
			return
		}

		session.SetSession(w, UserInfo.Username, UserInfo.Email, UserInfo.Password)

		w.Header().Set("Content-Type", "Application/JSON")

		err = email.SendMail(UserInfo.Email, UserInfo.Username)
		if err != nil {
			errorcheck.Nigger("File:register.go Error Sending mail :", err)
			fmt.Fprintf(w, `{"ok":false}`)
			return
		}
		fmt.Fprintf(w, `{"ok":true}`)

	}
}
