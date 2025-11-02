package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Keystring string = os.Getenv("JWT_KEY")
var Key []byte = []byte(Keystring)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Setjwtkey(w http.ResponseWriter, username string) {

	expat := time.Now().UTC().Add(1 * time.Hour)
	claims := Claims{}
	claims.Username = username
	claims.ExpiresAt = jwt.NewNumericDate(expat)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString(Key)
	if err != nil {
		fmt.Println("Errror creating token in jwt.go", err)
	}
	// fmt.Fprintf(w, tk)
	// fmt.Println(VerifyJWT(tk))

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tk,
		Path:     "/",
		Expires:  expat,
		HttpOnly: true,
		Secure:   false,                // set true only if you are using HTTPS
		SameSite: http.SameSiteLaxMode, // ✅ safer and more consistent than None without HTTPS
	})

	fmt.Println("✅ New JWT set for user:", username)

}

func VerifyJWT(tokenstring string) string {
	t, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (any, error) {
		return Key, nil
	})
	if err != nil {
		fmt.Println(err)
	}
	if !t.Valid {
		fmt.Println("Invalid")
	}
	c, _ := t.Claims.(*Claims)
	fmt.Println("Username from jwt:", c.Username)
	return c.Username

	// w.Header().Set("content-type", "application/json")
	// fmt.Fprintf(w, `{"key":%s,"username":%s}`, tokenstring, c.Username)
}
