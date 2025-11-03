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
			return
		}

		sessionid := cookie.Value
		fmt.Println("Cookie: ", cookie.Value)
		conn, err := db.Connect()
		fmt.Println(sessionid)
		if err != nil {
			fmt.Println("Error From sessionchecker [ Error Connecting to DB ]", err)

		}
		var expiry time.Time
		conn.QueryRow(context.Background(), "select expires from sessions where sessionid=$1", sessionid).Scan(&expiry)
		w.Header().Set("Content-Type", "application/json")
		if expiry.After(time.Now().UTC()) {
			cook.Active = true
			json.NewEncoder(w).Encode(cook)
		} else {
			cokkie := &http.Cookie{
				Name:     "sessionid",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			}
			http.SetCookie(w, cokkie)
			cook.Active = false
			json.NewEncoder(w).Encode(cook)
		}
	}

}

func CheckForVerifyPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	cook := Cookie_check{}

	cookie, err := r.Cookie("sessionid")
	if err != nil || cookie.Value == "" {
		fmt.Println("No valid session cookie found:", err)
		cook.Active = false
		_ = json.NewEncoder(w).Encode(cook)
		return
	}

	conn, err := db.Connect()
	if err != nil {
		fmt.Println("DB connection error:", err)
		cook.Active = false
		_ = json.NewEncoder(w).Encode(cook)
		return
	}
	defer conn.Close(context.Background())

	var expiry time.Time
	err = conn.QueryRow(context.Background(),
		"SELECT expires FROM sessions WHERE sessionid=$1",
		cookie.Value,
	).Scan(&expiry)

	if err != nil {
		fmt.Println("Session not found:", err)
		cook.Active = false
		_ = json.NewEncoder(w).Encode(cook)
		return
	}

	if expiry.After(time.Now().UTC()) {
		cook.Active = true
	} else {
		fmt.Println("Session expired, clearing cookie")
		http.SetCookie(w, &http.Cookie{
			Name:     "sessionid",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		cook.Active = false
	}

	_ = json.NewEncoder(w).Encode(cook)
}
