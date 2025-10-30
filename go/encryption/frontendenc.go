package encryption

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"vaultx/errorcheck"
	auth "vaultx/loginv1"
)

func Frontend_Enc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method no POST")
		return
	}
	cookie, err := r.Cookie("token")
	errorcheck.PrintError("Error getting cookie in Backend_Encryption.go", err)
	tokenstring := cookie.Value

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error getting file in frontend_enc() frontendenc.go", err)
	}
	filename := header.Filename
	username := auth.VerifyJWT(tokenstring)
	uploadpath := "/home/a4bhi/rawfiles/" + username
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
	fmt.Println(base64.StdEncoding.DecodeString(r.FormValue("file_iv")))

}
