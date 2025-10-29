package encryption

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"vaultx/db"
	"vaultx/errorcheck"
	auth "vaultx/loginv1"
)

func Backend_Encryption(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		fmt.Println("Requested method is not POST")
		return
	}

	cookie, err := r.Cookie("token")
	errorcheck.PrintError("Error getting cookie in Backend_Encryption.go", err)
	tokenstring := cookie.Value
	username := auth.VerifyJWT(tokenstring) //username from jwt
	mailid := GetMailidFromUsername(username)

	r.ParseMultipartForm(200 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error in r.FormFile()", err)
	}
	filename := header.Filename

	uploadpath := "/home/a4bhi/rawfiles" + username

	os.MkdirAll(uploadpath, os.ModePerm)

	dstpath := filepath.Join(uploadpath, filename)
	dst, err := os.Create(dstpath)
	if err != nil {
		fmt.Println("Error creating dstpath", err)
	}

	_, err = io.Copy(dst, file)

	if err != nil {
		fmt.Println("Error copying file.", err)
	}

}

func AesEnc(rawfilepath string) {
	rawfile, err := os.ReadFile(rawfilepath)
	errorcheck.PrintError("Error reading file in Aesaenc() in backendenc.go file", err)

}

func GetMailidFromUsername(username string) string {
	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in backendenc.go file ", err)
	var mailid string
	conn.QueryRow(context.TODO(), "select mailid from users where username=$1", username).Scan(&mailid)
	return mailid
}
