package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	username string
	jwt.RegisteredClaims
}

func Setjwtkey(w http.ResponseWriter, r *http.Request, username string) {
	var key string = os.Getenv("JWT_KEY")
	expat := time.Now().UTC().Add(1 * time.Hour)
	claims := Claims{}
	claims.username = username
	claims.ExpiresAt = expat
}
