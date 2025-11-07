package deleteacc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"vaultx/db"
	"vaultx/errorcheck"
	auth "vaultx/loginv1"
	"vaultx/masterkeys"

	"golang.org/x/crypto/bcrypt"
)

type AccDeleteResponse struct {
	Valid  bool   `json:"Valid"`
	Ok     bool   `json:"Ok"`
	Status string `json:"Status"`
}

func DeleteAcc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Println("Error expected method was DELETE but got : ", r.Method)
		return
	}
	var mailid string

	w.Header().Set("Content-Type", "application/json")

	ADR := AccDeleteResponse{}
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Error cant get cookie named token in DeleteAcc():", err)
		masterkeys.DeleteMasterkey(mailid)
		ADR.Valid = false
		ADR.Ok = false
		ADR.Status = ""
		json.NewEncoder(w).Encode(&ADR)
		return
	}

	tokenstring := cookie.Value

	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Error Invalid JWT in DeleteAcc():", err)
		masterkeys.DeleteMasterkey(mailid)
		ADR.Valid = false
		ADR.Ok = false
		ADR.Status = ""
		json.NewEncoder(w).Encode(&ADR)
		return
	}

	password := r.FormValue("password")
	var hashedpass []byte

	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in Deleteacc()", err)
	defer conn.Close(context.TODO())

	conn.QueryRow(context.TODO(), "select hashed_passed , mailid from users where username=$1", username).Scan(&hashedpass, &mailid)

	err = bcrypt.CompareHashAndPassword(hashedpass, []byte(password))
	if err != nil {
		fmt.Println("Incorrect Password in DeleteAcc()", err)
		ADR.Valid = true
		ADR.Ok = false
		ADR.Status = ""
		json.NewEncoder(w).Encode(&ADR)
		return

	}

	_, err = conn.Exec(context.TODO(), "Delete from users where username = $1", username)
	if err != nil {
		fmt.Println("Error Deleting acc from users table in DeleteAcc()", err)
		ADR.Valid = false
		ADR.Ok = false
		ADR.Status = ""
		json.NewEncoder(w).Encode(&ADR)
		return

	}

	_, err = conn.Exec(context.TODO(), "Delete from files where username=$1", username)
	if err != nil {
		fmt.Println("Error Deleting acc from users table in DeleteAcc()", err)
		ADR.Valid = false
		ADR.Ok = false
		ADR.Status = ""
		json.NewEncoder(w).Encode(&ADR)
		return

	}

	path := "/home/a4bhi/encrypted_files/" + username

	os.RemoveAll(path)

	clearCookie(w)
	masterkeys.DeleteMasterkey(mailid)

	ADR.Valid = true
	ADR.Ok = true
	ADR.Status = "Deleted"

	json.NewEncoder(w).Encode(&ADR)

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
