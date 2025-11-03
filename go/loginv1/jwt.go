package auth

import (
	"encoding/json"
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
type OnLoadResponse struct {
	Valid    bool   `json:"Valid"`
	Username string `json:"Username"`
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

func VerifyJWT(tokenstring string) (string, bool) {
	t, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (any, error) {
		return Key, nil
	})
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	if !t.Valid {
		fmt.Println("Invalid")
		return "", false
	}
	c, _ := t.Claims.(*Claims)
	fmt.Println("Username from jwt:", c.Username)
	return c.Username, true

	// w.Header().Set("content-type", "application/json")
	// fmt.Fprintf(w, `{"key":%s,"username":%s}`, tokenstring, c.Username)
}

func CheckJWT_ONLOAD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Expected Method Was GET but GOT : ", r.Method, " In CheckJWT_ONLOAD()")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	check := OnLoadResponse{}
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Error No Cookie Found in CheckJWT_ONLOAD():", err)
		check.Valid = false
		json.NewEncoder(w).Encode(&check)
		return
	}

	tokenstring := cookie.Value

	username, ok := VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Error Invalid JWT TOKEN  in CheckJWT_ONLOAD():", err)
		check.Valid = false
		json.NewEncoder(w).Encode(&check)
		return
	}

	check.Valid = true
	check.Username = username
	json.NewEncoder(w).Encode(&check)

}
