package decryption

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"vaultx/db"
	"vaultx/encryption"
	"vaultx/errorcheck"
	auth "vaultx/loginv1"
	"vaultx/masterkeys"
)

type FileDetails struct {
	FileID   int    `json:"FileID"`
	FileName string `json:"FileName"`
}

func Backend_Decryption(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method is GET expected was POST in backend_dec() ")
		return
	}
	fd := FileDetails{}
	// body , _ := io.ReadAll(r.Body)
	json.NewDecoder(r.Body).Decode(&fd)
	cookie, err := r.Cookie("token")
	errorcheck.PrintError("Error getting jwt cookie in Backend_Decryption() backenddec.go: ", err)
	tokenstring := cookie.Value

	username := auth.VerifyJWT(tokenstring)
	mailid := encryption.GetMailidFromUsername(username)

	masterkey := masterkeys.LoadMasterKey(mailid)

	ciphertext, err, file_iv, enc_file_key, enc_file_key_iv := GetFileFromDB(fd.FileID, fd.FileName, username)
	if err != nil {
		fmt.Println("Error from GetFileFromDb() backenddec.go", err)
		return
	}

	//#####FILE KEY DECRRYPTION############
	block, err := aes.NewCipher(masterkey)
	if err != nil {
		fmt.Println("Error creating block backend decryption:", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating gcm backend decryption:", err)
	}

	filekey, err := gcm.Open(nil, enc_file_key_iv, enc_file_key, nil)
	if err != nil {
		fmt.Println("Error decrypting filekey backend decryption: ", err)
	}

	//##########################################
	AesDec(ciphertext, filekey, file_iv, w, fd.FileName)

}

func GetFileFromDB(Fileid int, FileName string, Username string) (ciphertext []byte, err error, fileiv []byte, encfilekey []byte, encfilekey_iv []byte) {
	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in FetFileFromDB() backenddec.go: ", err)
	defer conn.Close(context.TODO())
	var filepath string
	var file_iv, enc_file_key, enc_file_key_iv string
	conn.QueryRow(context.TODO(), "select filepath,file_iv,enc_file_key,enc_file_key_iv from files where file_id=$1 and filename=$2 and username=$3", Fileid, FileName, Username).Scan(&filepath, &file_iv, &enc_file_key, &enc_file_key_iv)
	file, err := os.Open(filepath)
	errorcheck.PrintError("Error Opening file in GetFileFromDB() backenddec.go: ", err)
	ciphertext, err = io.ReadAll(file)
	errorcheck.PrintError("Error Reading from File in GetFileFromDB() backendec.go: ", err)
	fileiv, err = hex.DecodeString(file_iv)
	encfilekey, err = hex.DecodeString(enc_file_key)
	encfilekey_iv, err = hex.DecodeString(enc_file_key_iv)

	return ciphertext, err, fileiv, encfilekey, encfilekey_iv

}

func AesDec(ciphertext []byte, filekey []byte, fileiv []byte, w http.ResponseWriter, filename string) {
	block, err := aes.NewCipher(filekey)
	if err != nil {
		fmt.Println("Error creating block:", err)
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating gcm:", err)
		return
	}

	finalfile, err := gcm.Open(nil, fileiv, ciphertext, nil)
	if err != nil {
		fmt.Println("Error decrypting file:", err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(finalfile)))
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	_, err = io.Copy(w, bytes.NewReader(finalfile))
	if err != nil {
		fmt.Println("Error sending file in AesDec() backenddec.go : ", err)
		return
	}
}
