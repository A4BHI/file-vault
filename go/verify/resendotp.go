package verify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"vaultx/email"
)

type CountnTime struct {
	Count      int
	LockedTime time.Time
}
type ResponseToJs struct {
	Limit_Reached bool  `json:"limit-reached"`
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
		mailid := getSession(r)
		namme := strings.Split(mailid, "@")
		c, ok := GetCount(mailid)
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
			unloacktime := time.Now().Add(30 * time.Minute)
			StoreCount(mailid, count, unloacktime)
			res.Email_Send = false
			res.Limit_Reached = true
			res.Unlock_At = unloacktime.Unix()
			json.NewEncoder(w).Encode(res)
			return

		}

		res.Limit_Reached = false
		StoreCount(mailid, count, time.Time{})
		err := email.SendMail(mailid, namme[0])
		if err != nil {
			fmt.Println("Failed to send email:", err)
			res.Email_Send = false

		} else {

			res.Email_Send = true

		}

		json.NewEncoder(w).Encode(res)

	}
}
