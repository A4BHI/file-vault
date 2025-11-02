package showfiles

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"vaultx/db"
	auth "vaultx/loginv1"
)

type FileID struct {
	ID int `json:"id"`
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Println("Expected Request Was DELETE BUT GOT :", r.Method)
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Error getting cookie in DeleteFile()", err)
	}
	tokenstring := cookie.Value

	fileid := FileID{}
	json.NewDecoder(r.Body).Decode(&fileid)

	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Wrong JWT in DeleteFile() deletefile.go")
		return
	}
	filepath := GetFilePath(fileid.ID, username)
	fmt.Println(fileid.ID)
	fmt.Println(filepath)
	err = os.Remove(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File already deleted, ignoring.")
		} else {
			fmt.Println("Error deleting file:", err)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted successfully",
	})

}

func GetFilePath(id int, username string) string {
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Error connecting to db in GetFilePath() delerefile.go ", err)
	}

	var filepath string

	conn.QueryRow(context.TODO(), "select filepath from files where file_id = $1 and username=$2", id, username).Scan(&filepath)
	return filepath

}
