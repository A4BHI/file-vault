package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"vaultx/db"
	"vaultx/errorcheck"
	auth "vaultx/loginv1"
	"vaultx/masterkeys"
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

	AesEnc(dstpath, outputpath, mailid, filename, username)

}

func AesEnc(rawfilepath string, outputpath string, mailid string, filename string, username string) {
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

	file_iv := hex.EncodeToString(iv)

	enc_filekey_hex, enc_filekey_iv := Encrypt_Filekey(mailid, filekey, iv)
	StoreToFilesTable(file_iv, enc_filekey_hex, enc_filekey_iv, filename, encfilepath, username)

}

func GetMailidFromUsername(username string) string {
	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in backendenc.go file ", err)
	var mailid string
	conn.QueryRow(context.TODO(), "select mailid from users where username=$1", username).Scan(&mailid)
	return mailid
}

func Encrypt_Filekey(mailid string, filekey []byte, fileiv []byte) (string, string) {
	masterkey := masterkeys.LoadMasterKey(mailid)

	block, err := aes.NewCipher(masterkey)
	errorcheck.PrintError("Error creating block for encrypting fiekey", err)

	gcm, err := cipher.NewGCM(block)
	errorcheck.PrintError("Error getting gcm in Encrypt_Filekey() in backendenc.go : ", err)

	filekeyiv := make([]byte, 32)
	rand.Read(filekeyiv)

	encrypted_filekey := gcm.Seal(nil, filekeyiv, filekey, nil)
	fmt.Println(encrypted_filekey)

	enc_filekey_hex := hex.EncodeToString(encrypted_filekey)
	enc_filekey_iv := hex.EncodeToString(filekeyiv)

	return enc_filekey_hex, enc_filekey_iv

	// fmt.Println("MasterKey : ", masterkey)

}

func StoreToFilesTable(file_iv string, enc_filekey_hex string, enc_filekey_iv string, filename string, filepath string, username string) {

	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in backendenc.go Storetofilestable(): ", err)
	Uploaded_At := time.Now().UTC()
	var ownerID int
	conn.QueryRow(context.TODO(), "Select userid from users where username=$1", username).Scan(&ownerID)

	_, err = conn.Exec(
		context.TODO(),
		`INSERT INTO files 
        (Owner_ID, Username, FileName, FilePath, File_IV, Enc_File_Key, Enc_File_Key_IV,Uploaded_At) 
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		ownerID, username, filename, filepath, file_iv, enc_filekey_hex, enc_filekey_iv, Uploaded_At,
	)

	errorcheck.PrintError("Error adding details to files table in StoreToFilesTable() function backendenc.go : ", err)

}
