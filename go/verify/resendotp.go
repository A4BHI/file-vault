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

func StoreCount(Email string, Count int) {
	ResendLimit.Store(Email, Count)
}

func GetCount(Email string) (any, bool) {
	return ResendLimit.Load(Email)
}

func ResendOtp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("content-type", "application/json")
		mailid := getSession(r)
		namme := strings.Split(mailid, "@")
		c, ok := GetCount(mailid)
		count := 0
		if ok {
			count = c.(int)
		}
		res := ResponseToJs{}

		if count >= 3 {
			res.Limit_Reached = true
			res.Email_Send = false
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Limit_Reached = false
		count += 1
		StoreCount(mailid, count)

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
