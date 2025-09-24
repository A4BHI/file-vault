package otps

import (
	"sync"
)

var m sync.Map

func StoreOtp(key string, value string) {
	m.Store(key, value)
}

func GetOtp(key string) any {
	otp, _ := m.Load(key)
	return otp
}

func DeleteOtp(key string) {
	m.Delete(key)
}

func CompareOtp(key string, otp string) bool {
	o, _ := m.Load(key)

	if o == otp {
		return true
	}

	return false
}
