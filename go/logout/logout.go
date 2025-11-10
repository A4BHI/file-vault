package logout

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vaultx/db"
	auth "vaultx/loginv1"
	"vaultx/masterkeys"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Expected Method Was Post But Got:", r.Method)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Cant access cookie in Logout() logout.go : ", err)
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}

	tokenstring := cookie.Value

	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Invalid Cookie in Logout() logout.go [Redirecting to index.html]")
		clearCookie(w)
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return

	}
	fmt.Println("Username from Logout For Debug:", username)
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("error in logout() cannot connect to db:", err)
		return
	}

	var mailid string
	conn.QueryRow(context.TODO(), "select mailid from users where username=$1", username).Scan(&mailid)
	masterkeys.DeleteMasterkey(mailid)
	clearCookie(w)

	http.Redirect(w, r, "/index.html", http.StatusSeeOther)

}

func clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
}
