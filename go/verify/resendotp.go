package verify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"vaultx/email"
	"vaultx/session"
)

type CountnTime struct {
	Count      int
	LockedTime time.Time
}
type ResponseToJs struct {
	Limit_Reached bool  `json:"limit_reached"`
	Email_Send    bool  `json:"email_send"`
	Unlock_At     int64 `json:"unlock_at"`
}

var ResendLimit sync.Map

func StoreCount(Email string, Count int, LockTime time.Time) {
	ResendLimit.Store(Email, CountnTime{Count: Count, LockedTime: LockTime})
}

func GetCount(Email string) (CountnTime, bool) {
	r, ok := ResendLimit.Load(Email)

	if !ok {
		return CountnTime{}, false
	} else {
		return r.(CountnTime), true
	}
}

func ResendOtp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		sessioid := getSession(r)
		username, mail, _, _, _, ok := session.GetSession(sessioid)

		if !ok {
			fmt.Println("Session error in ResendOtp()")
			return
		}

		c, ok := GetCount(mail)
		count := 0
		if ok {
			count = c.Count
		}
		res := ResponseToJs{}

		if count >= 4 {
			if time.Now().Before(c.LockedTime) {
				res.Limit_Reached = true
				res.Email_Send = false
				res.Unlock_At = c.LockedTime.Unix()
				json.NewEncoder(w).Encode(res)
				return
			} else {
				count = 0
			}

		}

		count++
		if count >= 4 {
			unlocktime := time.Now().Add(30 * time.Minute)
			StoreCount(mail, count, unlocktime)
			res.Email_Send = false
			res.Limit_Reached = true
			res.Unlock_At = unlocktime.Unix()
			json.NewEncoder(w).Encode(res)
			return

		}

		res.Limit_Reached = false
		StoreCount(mail, count, time.Time{})
		err := email.SendMail(mail, username)
		if err != nil {
			fmt.Println("Failed to send email:", err)
			res.Email_Send = false

		} else {

			res.Email_Send = true

		}

		json.NewEncoder(w).Encode(res)

	}
}
