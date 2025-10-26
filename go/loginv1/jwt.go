package auth

import (
	"fmt"
	"net/http"
	"os"
)

var key string = os.Getenv("JWT_KEY")

func Setjwtkey(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, os.Getenv("JWT_KEY"))
	fmt.Print(key)
	fmt.Print("test")
	fmt.Fprint(w, "test")

}
