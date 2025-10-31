package decryption

import (
	"context"
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
	json.NewDecoder(r.Body).Decode(fd)
	cookie, err := r.Cookie("token")
	errorcheck.PrintError("Error getting jwt cookie in Backend_Decryption() backenddec.go: ", err)
	tokenstring := cookie.Value

	username := auth.VerifyJWT(tokenstring)
	mailid := encryption.GetMailidFromUsername(username)

	masterkey := masterkeys.LoadMasterKey(mailid)

	ciphertext, err := GetFileFromDB(fd.FileID, fd.FileName, username)

}

func GetFileFromDB(Fileid int, FileName string, Username string) (ciphertext []byte, err error) {
	conn, err := db.Connect()
	errorcheck.PrintError("Error connecting to db in FetFileFromDB() backenddec.go: ", err)
	defer conn.Close(context.TODO())
	var filepath string
	conn.QueryRow(context.TODO(), "select filepath from files where file_id=$1 and filename=$2 and username=$3", Fileid, FileName, Username).Scan(&filepath)
	file, err := os.Open(filepath)
	errorcheck.PrintError("Error Opening file in GetFileFromDB() backenddec.go: ", err)
	ciphertext, err = io.ReadAll(file)
	errorcheck.PrintError("Error Reading from File in GetFileFromDB() backendec.go: ", err)

	return ciphertext, err

}
