package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Creds struct {
	email    string `json:"mailid"`
	password string `json:"password"`
}

func login(w http.Response, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error from Login.go login function: ", err)
		}

		creds := Creds{}

		json.Unmarshal(body, &creds)

	}

}
