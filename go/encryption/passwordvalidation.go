package encryption

import (
	"fmt"
	"net/http"
)

func ValidatePass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Error Expected Method was POST BUT GOT GET")
		return
	}
}
