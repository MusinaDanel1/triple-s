package bucket

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func deleteBucket(bucketName string, csvFilePath string, dataDir string) error {
	// Проверяем, существует ли ведро
	bucketPath := fmt.Sprintf("%s/%s", dataDir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return fmt.Errorf("bucket not found")
	}

	// Проверяем, есть ли в ведре объекты
	files, err := os.ReadDir(bucketPath)
	if err != nil {
		return fmt.Errorf("error reading bucket directory: %v", err)
	}
	if len(files) > 0 {
		return fmt.Errorf("bucket is not empty, delete objects before deleting the bucket")
	}

	// Удаляем директорию ведра
	err = os.RemoveAll(bucketPath)
	if err != nil {
		return fmt.Errorf("error deleting bucket directory: %v", err)
	}

	// Читаем существующие ведра из CSV
	file, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	// Записываем обратно в CSV все ведра, кроме удаляемого
	var updatedRecords [][]string
	for _, record := range records {
		if record[0] != bucketName {
			updatedRecords = append(updatedRecords, record)
		}
	}

	// Записываем обновленные записи в CSV файл
	file, err = os.OpenFile(csvFilePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("error opening CSV file for writing: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(updatedRecords)
	if err != nil {
		return fmt.Errorf("error writing to CSV file: %v", err)
	}

	return nil
}

// DeleteBucketHandler обрабатывает HTTP-запросы на удаление ведра.
func DeleteBucketHandler(w http.ResponseWriter, r *http.Request, dataDir, bucketName string) {
	// bucketName := r.URL.Path[len("/delete/"):] // Измените на len("/buckets/delete/") для более точного извлечения имени

	if bucketName == "" {
		http.Error(w, "400 Bad Request: Missing bucket name in the URL", http.StatusBadRequest)
		return
	}

	csvFilePath := filepath.Join(dataDir, "buckets.csv")

	err := deleteBucket(bucketName, csvFilePath, dataDir)
	if err != nil {
		if err.Error() == "bucket not found" {
			http.Error(w, "404 Not Found: Bucket not found", http.StatusNotFound)
		} else if strings.Contains(err.Error(), "bucket is not empty") {
			http.Error(w, "400 Bad Request: "+err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
