package session

import (
	"context"
	"fmt"
	"time"
	"vaultx/db"
	"vaultx/otps"
)

func CleanSessionFromDB() {
	var sessionids []string

	query := "select sessionid,expires,mailid from sessions"
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Error in CleansessionfromDB():", err)
	}

	row, err := conn.Query(context.TODO(), query)
	if err != nil {
		fmt.Println("Error conecting to db in CleanSessionFromDB()", err)
	}
	for row.Next() {
		var sessionid string
		var expires time.Time
		var mail string
		row.Scan(&sessionid, &expires, &mail)
		if time.Now().After(expires) {
			sessionids = append(sessionids, sessionid)
			otps.DeleteOtp(mail)

		}
	}

	for i := range sessionids {
		conn.Exec(context.Background(), "delete from sessions where sessionid=$1", sessionids[i])
		fmt.Printf("Session %s has been deleted \n", sessionids[i])
	}
}
