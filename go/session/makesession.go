package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"vaultx/db"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func SetSession(w http.ResponseWriter, username string, mailid string, password string) bool {
	randomb := make([]byte, 32)
	rand.Read(randomb)
	sessionid := hex.EncodeToString(randomb)
	createdat := time.Now().UTC()
	expires := createdat.Add(30 * time.Minute)
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		fmt.Println("makesession.go setSessionFunction salt:", err)
	}
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("makesession.go setSessionFunction hashedpass:", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionid",
		Value:    sessionid,
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
	})

	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Error connecting to db in makesession.go", err)
		return false
	}

	_, err = conn.Exec(context.TODO(), "insert into sessions  values($1,$2,$3,$4,$5,$6,$7)", username, mailid, sessionid, hashedpass, salt, createdat, expires)
	if err != nil {
		fmt.Println("Error executing query for adding session to db in makesession.go", err)
		return false
	}
	return true
}

func DeleteSession(sessionid string) bool {
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Error connecting to db from makesession.go in DeleteSession func :", err)
		return false
	}

	_, err = conn.Exec(context.TODO(), "DELETE from sessions where sessionid=$1", sessionid)
	if err != nil {
		fmt.Println("Error deleting session from db file:makesession.go:", err)
		return false
	}

	return true

}

func GetSession(sessionid string) (user string, mail string, hashedpas string, upp []byte, ok bool) {
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Error connecting to db in getsession():", err)
		return "", "", "", nil, false
	}
	var salt []byte
	var username, mailid, hashedpass string
	err = conn.QueryRow(context.TODO(), "select username,mailid,hashedpass,salt from sessions where sessionid = $1", sessionid).Scan(&username, &mailid, &hashedpass, &salt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", nil, false
		}
	}
	return username, mailid, hashedpass, salt, true
}
