package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"vaultx/db"

	"vaultx/errorcheck"
	auth "vaultx/loginv1"
	"vaultx/masterkeys"
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
	iv, err := base64.StdEncoding.DecodeString(r.FormValue("file_iv"))
	enc_file_key, err := base64.StdEncoding.DecodeString(r.FormValue("enc_file_key"))
	key_iv, err := base64.StdEncoding.DecodeString(r.FormValue("key_iv"))
	salt, err := base64.StdEncoding.DecodeString(r.FormValue("salt"))

	mailid := GetMailidFromUsername(username)

	masterkey := masterkeys.LoadMasterKey(mailid)

	// filepath := "/home/a4bhi/rawfiles/" + username + "/" + filename

	block, err := aes.NewCipher(masterkey)
	if err != nil {
		fmt.Println("Error creating block:", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating gcm:", err)
	}

	fiekey, err := gcm.Open(nil, key_iv, enc_file_key, nil)
	if err != nil {
		fmt.Println("Error decrypting filekey: ", err)
	}

	fmt.Println("Decrypted FILE KEY : ", fiekey)

	fmt.Println("IV FROM FRONTEND: ", iv)
	fmt.Println("ENC FILE KEY: ", enc_file_key)
	fmt.Println("KEY IV", key_iv)
	fmt.Println("SALT FROM FRONTEND", salt)

}

func GetSalt(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	errorcheck.PrintError("Error getting cookie in Backend_Encryption.go", err)
	tokenstring := cookie.Value
	username := auth.VerifyJWT(tokenstring)
	fmt.Println(username)
	var salt []byte
	conn, _ := db.Connect()
	conn.QueryRow(context.TODO(), "select salt from users where username=$1", username).Scan(&salt)
	fmt.Println("Salt :", salt)
	// salt, _ := hex.DecodeString(salthex)
	salltBase64 := base64.StdEncoding.EncodeToString(salt)
	// fmt.Println("SALT FROM BACKEND", salt)
	// fmt.Println("SaltBASE64 : ", salltBase64)
	w.Header().Set("content-type", "text/plain")
	fmt.Fprint(w, salltBase64)

}
