package encryption

import (
	"context"
	"fmt"
	"net/http"
	"vaultx/db"
	"vaultx/errorcheck"
	auth "vaultx/loginv1"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Error Expected Method was POST BUT GOT GET")
		return
	}
	w.Header().Set("content-type", "text/plain")

	password := r.FormValue("password")
	fmt.Println(password)

	cookie, err := r.Cookie("token")
	errorcheck.PrintError("Error Getting Cookie in ValidatePass() passwordvalidation.go", err)

	tokenstring := cookie.Value

	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Wrong JWT in ValidatePass() PasswordValidation.go")
		return
	}

	var hashedpass []byte
	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in ValidatePass() passwordvalidation.go", err)
	defer conn.Close(context.TODO())

	conn.QueryRow(context.TODO(), "select hashed_passed from users where username=$1", username).Scan(&hashedpass)

	err = bcrypt.CompareHashAndPassword(hashedpass, []byte(password))
	if err != nil {
		fmt.Println("Incorrect Password in ValidatePass() passwordvalidation.go", err)
		fmt.Fprintf(w, "Fail")
		return

	}
	fmt.Println("PAssword crct", password)
	fmt.Fprintf(w, "Pass")

}
