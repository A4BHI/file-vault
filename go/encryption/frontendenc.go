package encryption

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"vaultx/db"

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
	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Wrong JWT in Frontend_Enc() frontenc.go")
		return
	}
	uploadpath := "/home/a4bhi/encrypted_files/" + username
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
	iv, err := base64.StdEncoding.DecodeString(r.FormValue("file_iv"))
	errorcheck.PrintError("Error Decoding From BASE64 to BYTES in frontendenc.go {iv}", err)

	enc_file_key, err := base64.StdEncoding.DecodeString(r.FormValue("enc_file_key"))
	errorcheck.PrintError("Error Decoding From BASE64 to BYTES in frontendenc.go {enc_file_key}", err)

	key_iv, err := base64.StdEncoding.DecodeString(r.FormValue("key_iv"))
	errorcheck.PrintError("Error Decoding From BASE64 to BYTES in frontendenc.go {key_iv}", err)

	salt, err := base64.StdEncoding.DecodeString(r.FormValue("salt"))
	errorcheck.PrintError("Error Decoding From BASE64 to BYTES in frontendenc.go {salt}", err)

	fmt.Println("Salt from Frontend: ", salt)

	iv_hex := hex.EncodeToString(iv)
	enc_file_key_hex := hex.EncodeToString(enc_file_key)
	enc_file_key_iv_hex := hex.EncodeToString(key_iv)

	StoreToFilesTable(iv_hex, enc_file_key_hex, enc_file_key_iv_hex, filename, dstpath, username)

}

func GetSalt(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	errorcheck.PrintError("Error getting cookie in Backend_Encryption.go", err)
	tokenstring := cookie.Value
	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Wrong JWT in GetSalt() frontendenc.go")
		return
	}
	fmt.Println(username)
	var salt []byte
	conn, _ := db.Connect()
	defer conn.Close(context.TODO())
	conn.QueryRow(context.TODO(), "select salt from users where username=$1", username).Scan(&salt)
	fmt.Println("Salt :", salt)

	salltBase64 := base64.StdEncoding.EncodeToString(salt)

	w.Header().Set("content-type", "text/plain")
	fmt.Fprint(w, salltBase64)

}
