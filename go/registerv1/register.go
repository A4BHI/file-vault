package registerv1

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"net/http"

	"vaultx/db"
	"vaultx/email"
	"vaultx/errorcheck"
)

type Creds struct {
	Username string `json:"username"`
	Email    string `json:"mailid"`
	Password string `json:"password"`
}

var UserInfo Creds

func setSession(email string, w http.ResponseWriter) {
	emailbytes := []byte(email)

	hashedMail := hex.EncodeToString(emailbytes)

	http.SetCookie(w, &http.Cookie{
		Name:     "mail",
		Value:    hashedMail,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 30),
	})

}
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("++++++++++Register++++++++++")
		body, err := io.ReadAll(r.Body)
		errorcheck.Nigger("Error From register.go in register Function:", err)

		var exists bool

		err = json.Unmarshal(body, &UserInfo)
		errorcheck.Nigger("Error in register.go Register function json unmarshal", err)

		fmt.Println(UserInfo.Password)

		conn, err := db.Connect()
		errorcheck.Nigger("Error creating connection to psql (Register.go file in Register function)", err)

		row := conn.QueryRow(context.Background(), "select exists( select 1 from users where username = $1 or mailid=$2)", UserInfo.Username, UserInfo.Email)
		row.Scan(&exists)

		if exists {
			w.Header().Set("Content-Type", "Application/JSON")
			fmt.Fprintf(w, `{"ok":false}`)
			return
		}

		setSession(UserInfo.Email, w)

		w.Header().Set("Content-Type", "Application/JSON")
		fmt.Fprintf(w, `{"ok":true}`)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		go func(emailAddr, username string) {
			err := email.SendMail(emailAddr, username)
			if err != nil {
				fmt.Println("Failed to send email:", err)
			}
		}(UserInfo.Email, UserInfo.Username)
		errorcheck.Nigger("File:register.go Error Sending mail :", err)

	}
}
