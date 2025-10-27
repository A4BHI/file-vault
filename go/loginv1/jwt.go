package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	username string `json:"username"`
	jwt.RegisteredClaims
}

func Setjwtkey(w http.ResponseWriter, r *http.Request) {
	var keystring string = os.Getenv("JWT_KEY")
	key := []byte(keystring)
	expat := time.Now().UTC().Add(1 * time.Hour)
	claims := Claims{}
	claims.username = "TEST"
	claims.ExpiresAt = jwt.NewNumericDate(expat)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString(key)
	if err != nil {
		fmt.Println("Errror creating token in jwt.go", err)
	}
	fmt.Fprintf(w, tk)

}
