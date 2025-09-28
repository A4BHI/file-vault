package registerv1

import (
	"context"
	"crypto/rand"
	"fmt"
	math "math/rand"
	"time"
	"vaultx/db"
	"vaultx/errorcheck"

	"golang.org/x/crypto/bcrypt"
)

func SaveToDB() bool {
	// creds := Creds{}

	userid := math.Intn(9000) + 1000

	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	errorcheck.Nigger("Error in rand.read(salt) on SaveToDB function:", err)

	HashedPass, err := bcrypt.GenerateFromPassword([]byte(UserInfo.Password), 12)
	if err != nil {
		fmt.Println("Error generating hashed password in function savetodb():", err)
	}

	createdat := time.Now().UTC()

	conn, err := db.Connect()
	errorcheck.Nigger("Error connecting to db file addtodb.go function:SaveToDB:", err)

	_, err = conn.Exec(context.Background(), "insert into users values($1,$2,$3,$4,$5,$6)", userid, UserInfo.Username, UserInfo.Email, salt, HashedPass, createdat)
	if err != nil {
		fmt.Println("Error executing sql query in SavetoDb function:", err)
		return false
	} else {
		return true
	}
}
