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

type FileInfo struct {
	FileID   int    `json:"FileID"`
	FileName string `json:"FileName"`
	FileType string `json:"FileType"`
	FileSize int64  `json:"FileSize"`
}

func ShowFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get username from JWT
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
		return
	}
	tokenstring := cookie.Value
	username, ok := auth.VerifyJWT(tokenstring)
	if !ok {
		fmt.Println("Wrong JWT in ShowFiles() myfiles.go")
		return
	}
	fmt.Println("For Debug:", username)

	// Fetch all file details for this user
	files, err := GetFileDetails(username)
	if err != nil {
		http.Error(w, "Error fetching files", http.StatusInternalServerError)
		fmt.Println("Error in GetFileDetails:", err)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func GetFileDetails(username string) ([]FileInfo, error) {
	conn, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("db connect error: %v", err)
	}
	defer conn.Close(context.TODO())

	rows, err := conn.Query(context.TODO(),
		"SELECT file_id, filename, filepath FROM files WHERE username=$1", username)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var files []FileInfo

	for rows.Next() {
		var fileID int
		var fileName, filePath string
		if err := rows.Scan(&fileID, &fileName, &filePath); err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}

		// Get file info
		info, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Stat error for", filePath, ":", err)
			continue
		}

		// Detect MIME type
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Open error:", err)
			continue
		}
		defer file.Close()

		buffer := make([]byte, 512)
		_, _ = file.Read(buffer)
		fileType := http.DetectContentType(buffer)

		files = append(files, FileInfo{
			FileID:   fileID,
			FileName: fileName,
			FileType: fileType,
			FileSize: info.Size(),
		})
	}

	return files, nil
}
