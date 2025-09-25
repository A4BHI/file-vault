package registerv1

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	// mathrand "math/rand"
	"net/http"
	// "time"
	"vaultx/db"
	"vaultx/errorcheck"
	// "golang.org/x/crypto/bcrypt"
)

type Creds struct {
	Username string `json:"username"`
	Email    string `json:"mailid"`
	Password string `json:"password"`
}

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

		creds := Creds{}
		var exists bool

		err = json.Unmarshal(body, &creds)
		errorcheck.Nigger("Error in register.go Register function json unmarshal", err)

		fmt.Println(creds.Password)

		// hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)

		// errorcheck.Nigger("Error generating hash password (register.go file Register() function )", err)

		conn, err := db.Connect()
		errorcheck.Nigger("Error creating connection to psql (Register.go file in Register function)", err)

		// createdat := time.Now()
		// userid := mathrand.Intn(9000) + 1000
		// salt := make([]byte, 32)
		// rand.Read(salt)
		row := conn.QueryRow(context.Background(), "select exists( select 1 from users where username = $1 or mailid=$2)", creds.Username, creds.Email)
		row.Scan(&exists)

		if exists {
			w.Header().Set("Content-Type", "Application/JSON")
			fmt.Fprintf(w, "Fail")
		} else {
			setSession(creds.Email, w)
			w.Header().Set("Content-Type", "Application/JSON")
			fmt.Fprintf(w, "Pass")

			// err := email.SendMail(creds.Email, creds.Username)
			errorcheck.Nigger("File:register.go Error Sending mail :", err)

		}

	}
}
