package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"triple-s/pkg/bucket"
	"triple-s/pkg/object"
)

func SetupRoutes(dataDir string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		switch len(pathParts) {
		case 2:
			bucketName := pathParts[1]
			if r.Method == http.MethodPut {
				bucket.CreateBucketHandler(w, r, dataDir, bucketName)
			} else if r.Method == http.MethodDelete {
				bucket.DeleteBucketHandler(w, r, dataDir, bucketName)
			} else if r.Method == http.MethodGet {
				bucket.ListAllBucketsHandler(w, r, dataDir)
			} else {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		case 3:
			bucketName := pathParts[1]
			objectKey := pathParts[2]
			if r.Method == http.MethodPut {
				object.UploadObjectHandler(w, r, dataDir, bucketName, objectKey)
			} else if r.Method == http.MethodGet {
				object.RetrieveObjectHandler(w, r, dataDir, bucketName, objectKey)
			} else if r.Method == http.MethodDelete {
				object.DeleteObjectHandler(w, r, dataDir, bucketName, objectKey)
			} else {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Bad Request: Invalid URL format", http.StatusBadRequest)
		}
	})
	return mux
}

func ValidatePort(port string) (int, error) {
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return 0, fmt.Errorf("Invalid port number: %s. Port must be an integer between 1 and 65535", port)
	}
	return portNum, nil
}
