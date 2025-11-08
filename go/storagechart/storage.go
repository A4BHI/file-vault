package storagechart

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	auth "vaultx/loginv1"
)

const totalLimitBytes = 1 * 1024 * 1024 * 1024 // 1 GB

type StorageResponse struct {
	UsedBytes   int64   `json:"used_bytes"`
	UsedPercent float64 `json:"used_percent"`
	FreeBytes   int64   `json:"free_bytes"`
	FreePercent float64 `json:"free_percent"`
}

func Sendata(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Error cannot get cookie name token in Senddata() storage.go:", err)
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return

	}
	tokenstring := cookie.Value

	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Invalid JWT FOUND inn SENDATA() ")
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return

	}
	path := filepath.Join("/home/a4bhi/encrypted_files/", username)
	usedBytes, err := dirSize(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if usedBytes > totalLimitBytes {
		usedBytes = totalLimitBytes // clamp to 100%
	}

	usedPercent := (float64(usedBytes) / float64(totalLimitBytes)) * 100
	freeBytes := totalLimitBytes - usedBytes
	freePercent := 100 - usedPercent

	resp := StorageResponse{
		UsedBytes:   usedBytes,
		UsedPercent: usedPercent,
		FreeBytes:   freeBytes,
		FreePercent: freePercent,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

// dirSize recursively sums file sizes inside a directory
func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
