package session

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"vaultx/db"
)

func SetSession(w http.ResponseWriter, username string, mailid string) bool {
	randomb := make([]byte, 32)
	sessionid := hex.EncodeToString(randomb)
	createdat := time.Now().UTC()
	expires := createdat.Add(30 * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:     sessionid,
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

	_, err = conn.Exec(context.TODO(), "insert into sessions (username,email,sessionid,createdat,expires) values($1,$2,$3,$4,$5)", username, mailid, sessionid, createdat, expires)
	if err != nil {
		fmt.Println("Error executing query for adding session to db in makesession.go", err)
		return false
	}
	return true
}
