package session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vaultx/db"
)

type Cookie_check struct {
	Active bool `json:"Active"`
}

func Check(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cook := Cookie_check{}
		cookie, err := r.Cookie("sessionid")
		if err != nil {
			fmt.Println("error from sessioncheck cant read cookie:", err)
		}

		sessionid := cookie.Value
		conn, err := db.Connect()
		fmt.Println(sessionid)
		if err != nil {
			fmt.Println("Error From sessionchecker [ Error Connecting to DB ]", err)

		}
		var expiry time.Time
		conn.QueryRow(context.Background(), "select expiry_time from sessions where sessionid=$1", sessionid).Scan(&expiry)
		w.Header().Set("Content-Type", "application/json")
		if expiry.After(time.Now().UTC()) {
			cook.Active = true
			json.NewEncoder(w).Encode(cook)
		} else {

			cook.Active = false
			json.NewEncoder(w).Encode(cook)
		}
	}

}
