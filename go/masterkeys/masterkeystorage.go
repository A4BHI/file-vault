package masterkeys

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"vaultx/db"

	"golang.org/x/crypto/pbkdf2"
)

var master sync.Map

func StoreKey(email string, masterkey []byte) {
	master.Store(email, masterkey)
}

func GenerateKey(mailid string, password string) (masterkey []byte) {
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Error in masterkeystirage.go GenerateKey function", err)
	}
	var salt []byte
	_, ok := master.Load(mailid)
	if ok {
		fmt.Println("Master key already in memory of this user")
		fmt.Println(master.Load(mailid))
		return
	}
	conn.QueryRow(context.Background(), "select salt from users where mailid=$1", mailid).Scan(&salt)
	byteForm := []byte(password)
	masterkey = pbkdf2.Key(byteForm, salt, 200000, 32, sha256.New)
	master.Store(mailid, masterkey)

	return masterkey

}

func LoadMasterKey(mailid string) []byte {
	value, ok := master.Load(mailid)
	if !ok {
		fmt.Println("Error getting masterkey LoadMasterKey func")
	}

	masterkey, ok := value.([]byte)
	if !ok {
		fmt.Println("Cannot convert any to byte")
	}

	return masterkey

}
