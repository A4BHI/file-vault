package otps

import (
	"sync"
)

var m sync.Map

func StoreOtp(key string, value string) {
	m.Store(key, value)
}
