package registerv1

import (
	"crypto/rand"
	"vaultx/db"
	"vaultx/errorcheck"
)

func SaveToDB() {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	errorcheck.Nigger("Error in rand.read(salt) on SaveToDB function:", err)

	c := Creds{}

	conn, err := db.Connect()
	errorcheck.Nigger("Error connecting to db file addtodb.go function:SaveToDB:", err)

}
