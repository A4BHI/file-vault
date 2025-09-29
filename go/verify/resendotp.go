package verify

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"vaultx/email"
)

var ResendLimit sync.Map

var count int = 0

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
		c, _ := GetCount(mailid)
		if c.(int) <= 3 {
			count += 1
			StoreCount(mailid, count)
			fmt.Fprintf(w, `{"Status":"Sending"}`)
			go func(emailAddr, username string) {
				err := email.SendMail(emailAddr, username)
				if err != nil {
					fmt.Println("Failed to send email:", err)
				} else {
					fmt.Fprintf(w, `{"Email":"Send"}`)
				}
			}(mailid, namme[0])

		} else {
			fmt.Fprintf(w, `{"Status":"Limit-Reached"}`)
		}

	}
}
