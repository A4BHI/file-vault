package registerv1

import (
	"context"
	"fmt"
	math "math/rand"
	"net/http"
	"time"
	"vaultx/db"
	"vaultx/errorcheck"
	"vaultx/session"
)

func SaveToDB(r *http.Request) bool {
	// creds := Creds{}
	cookie, _ := r.Cookie("sessionid")
	sessionid := cookie.Value
	//store password hash and salt also in db
	username, mailid, hashedpass, salt, ok := session.GetSession(sessionid)
	if !ok {
		fmt.Println("Error calling getSession() in addtodb")
		return false
	}
	userid := math.Intn(9000) + 1000

	createdat := time.Now().UTC()

	conn, err := db.Connect()
	errorcheck.Nigger("Error connecting to db file addtodb.go function:SaveToDB:", err)

	_, err = conn.Exec(context.Background(), "insert into users values($1,$2,$3,$4,$5,$6)", userid, username, mailid, salt, hashedpass, createdat)
	if err != nil {
		fmt.Println("Error executing sql query in SavetoDb function:", err)
		return false
	} else {
		return true
	}
}
