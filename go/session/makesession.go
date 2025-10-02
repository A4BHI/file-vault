package session

import "encoding/hex"

func SetSession(username string, mailid string) bool {
	randomb := make([]byte, 32)
	sessionid := hex.EncodeToString(randomb)

}
