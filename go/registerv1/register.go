package registerv1

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	mathrand "math/rand"
	"net/http"
	"time"
	"vaultx/db"
	"vaultx/errorcheck"

	"golang.org/x/crypto/bcrypt"
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

		creds := Creds{}

		err = json.Unmarshal(body, &creds)
		errorcheck.Nigger("Error in register.go Register function json unmarshal", err)

		fmt.Println(creds.Password)

		hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)

		errorcheck.Nigger("Error generating hash password (register.go file Register() function )", err)

		conn, err := db.Connect()
		errorcheck.Nigger("Error creating connection to psql (Register.go file in Register function)", err)

		createdat := time.Now()
		userid := mathrand.Intn(9000) + 1000
		salt := make([]byte, 32)
		rand.Read(salt)
		_, err = conn.Exec(context.Background(), "insert into users values($1,$2,$3,$4,$5,$6)", userid, creds.Username, creds.Email, salt, hashed, createdat)
		if err == nil {
			w.Header().Set("Content-Type", "Apllication/JSON")
			fmt.Fprint(w, "success")
		}
	}
}
