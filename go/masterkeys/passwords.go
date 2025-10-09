package masterkeys

import "sync"

var Pass sync.Map

func StorePassword(mailid string, password string) {
	Pass.Store(mailid, password)
}

func DeletePassword(mailid string) {
	Pass.Delete(mailid)
}

func GetPassword(mailid string) (password any, ok bool) {
	return Pass.Load(mailid)
}
