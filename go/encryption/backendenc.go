package encryption

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	username := auth.VerifyJWT(tokenstring)
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error in r.FormFile()", err)
	}
	filename := header.Filename

	uploadpath := "/home/a4bhi/" + username

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

	if err != nil {
		fmt.Println("Error getting file in Backend_Encryption.go file")
		return
	}

}

func AesEnc() {

}
