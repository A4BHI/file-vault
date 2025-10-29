package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

	uploadpath := "/home/a4bhi/rawfiles/" + username

	outputpath := "/home/a4bhi/encrypted_files/" + username

	os.MkdirAll(uploadpath, os.ModePerm)

	os.MkdirAll(outputpath, os.ModePerm)

	dstpath := filepath.Join(uploadpath, filename)
	dst, err := os.Create(dstpath)
	if err != nil {
		fmt.Println("Error creating dstpath", err)
	}

	defer file.Close()
	defer dst.Close()

	_, err = io.Copy(dst, file)

	if err != nil {
		fmt.Println("Error copying file.", err)
	}

	AesEnc(dstpath, outputpath, mailid, filename)

}

func AesEnc(rawfilepath string, outputpath string, mailid string, filename string) {
	rawfile, err := os.ReadFile(rawfilepath)
	errorcheck.PrintError("Error reading file in Aesaenc() in backendenc.go file", err)

	filekey := make([]byte, 32)
	_, err = rand.Read(filekey)
	errorcheck.PrintError("Error creating file key in aesenc() backendenc.go:", err)

	block, err := aes.NewCipher(filekey)
	errorcheck.PrintError("error block in aesenc() backendenc.go:", err)

	gcm, err := cipher.NewGCM(block)
	errorcheck.PrintError("Error creating cipherblock(gcm) in aesenc() backendenc.go", err)

	iv := make([]byte, gcm.NonceSize())
	_, err = rand.Read(iv)
	errorcheck.PrintError("Error creating iv aesenc() backendenc.go", err)

	ciphertext := gcm.Seal(nil, iv, rawfile, nil)

	encfilepath := filepath.Join(outputpath, filename)

	outputfile, err := os.Create(encfilepath)
	errorcheck.PrintError("Error creating outputfile in aesenc() backendenc.go:", err)

	defer outputfile.Close()

	outputfile.Write(ciphertext)

	os.Remove(rawfilepath)

}

func GetMailidFromUsername(username string) string {
	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in backendenc.go file ", err)
	var mailid string
	conn.QueryRow(context.TODO(), "select mailid from users where username=$1", username).Scan(&mailid)
	return mailid
}
