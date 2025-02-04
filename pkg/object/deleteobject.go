package object

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// DeleteObjectHandler обрабатывает удаление объекта из бакета.
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request, bucketDir, bucketName, objectKey string) {
	if bucketName == "" || objectKey == "" {
		http.Error(w, "400 Bad Request: Missing bucket name or object key", http.StatusBadRequest)
		return
	}

	// 2. Проверка существования ведра
	bucketPath := filepath.Join(bucketDir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found: Bucket does not exist", http.StatusNotFound)
		return
	}

	// 3. Проверка существования объекта
	objectPath := filepath.Join(bucketPath, objectKey)
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found: Object does not exist", http.StatusNotFound)
		return
	}

	// 4. Удаление объекта
	if err := os.Remove(objectPath); err != nil {
		http.Error(w, "500 Internal Server Error: Unable to delete object", http.StatusInternalServerError)
		return
	}

	// 5. Обновление метаданных в CSVobjectCSVPath := fmt.Sprintf("data/%s/objects.csv", bucketName)
	if err := deleteObjectMetadata(bucketDir, bucketName, objectKey); err != nil {
		http.Error(w, "500 Internal Server Error: Unable to update metadata", http.StatusInternalServerError)
	}

	// 6. Возвращаем успешный ответ
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// deleteObjectMetadata обновляет CSV файл, удаляя запись о объекте
func deleteObjectMetadata(bucketDir, bucketName, objectKey string) error {
	objectCSVPath := fmt.Sprintf("%s/%s/objects.csv", bucketDir, bucketName)

	// Открываем файл метаданных объектов для чтения
	file, err := os.Open(objectCSVPath)
	if err != nil {
		return fmt.Errorf("unable to open objects metadata file: %v", err)
	}
	defer file.Close()

	// Читаем все строки из CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read objects metadata: %v", err)
	}

	// Записываем обратно в CSV все ведра, кроме удаляемого
	var updatedRecords [][]string
	for _, record := range records {
		if record[0] != objectKey {
			updatedRecords = append(updatedRecords, record)
		}
	}

	file, err = os.OpenFile(objectCSVPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0o644)
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
