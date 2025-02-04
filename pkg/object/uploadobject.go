package object

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

type ObjectMetadata struct {
	XMLName      xml.Name `xml:"Object"`
	Key          string   `xml:"Key"`
	Size         int64    `xml:"Size"`
	LastModified string   `xml:"LastModified"`
	ContentType  string   `xml:"ContentType"`
}

// UploadObjectHandler обрабатывает загрузку объекта в ведро
func UploadObjectHandler(w http.ResponseWriter, r *http.Request, bucketDir, bucketName, objectKey string) {
	// 1. Проверка имени ведра и ключа объекта
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

	// 3. Валидация ключа объекта
	if err := validateObjectKey(objectKey); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Сохранение объекта в файловую систему (перезаписывается, если уже существует)
	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Create(objectPath)
	if err != nil {
		http.Error(w, "500 Internal Server Error: Unable to create object file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 5. Запись данных объекта
	if _, err := io.Copy(file, r.Body); err != nil {
		http.Error(w, "500 Internal Server Error: Unable to write object data", http.StatusInternalServerError)
		return
	}

	// 6. Получаем информацию о объекте
	objectInfo, err := os.Stat(objectPath)
	if err != nil {
		http.Error(w, "500 Internal Server Error: Unable to get object info", http.StatusInternalServerError)
		return
	}

	// 7. Обновление метаданных объекта в CSV
	err = updateObjectMetadata(bucketName, objectKey, objectInfo.Size(), bucketDir, r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "500 Internal Server Error: Unable to update object metadata", http.StatusInternalServerError)
		return
	}

	// 8. Возвращаем успешный ответ
	objectMetadata := ObjectMetadata{
		Key:          objectKey,
		Size:         objectInfo.Size(),
		LastModified: objectInfo.ModTime().Format(time.RFC3339),
		ContentType:  r.Header.Get("Content-Type"),
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(objectMetadata); err != nil {
		http.Error(w, "500 Internal Server Error: Unable to encode XML", http.StatusInternalServerError)
	}
}

// validateObjectKey проверяет, соответствует ли ключ объекта правилам
func validateObjectKey(objectKey string) error {
	keyPattern := `^[a-zA-Z0-9._-]{1,255}$`
	if matched, _ := regexp.MatchString(keyPattern, objectKey); !matched {
		return fmt.Errorf("400 Bad Request: Object key must be 1-255 characters long and can only contain letters, numbers, underscores, hyphens, and periods")
	}
	return nil
}

// updateObjectMetadata обновляет CSV файл с метаданными объектов
func updateObjectMetadata(bucketName, objectKey string, size int64, dataDir, contentType string) error {
	objectCSVPath := fmt.Sprintf("%s/%s/objects.csv", dataDir, bucketName)

	// Если Content-Type не был передан, определяем его по расширению файла
	if contentType == "" {
		ext := path.Ext(objectKey)
		contentType = mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream" // По умолчанию
		}
	}

	// Чтение существующих записей из CSV
	records, err := readCSV(objectCSVPath)
	if err != nil {
		return err
	}

	// Проверяем, существует ли запись с таким объектом, и удаляем её
	var updatedRecords [][]string
	for _, record := range records {
		if record[0] != objectKey { // Убираем запись с удаляемым объектом
			updatedRecords = append(updatedRecords, record)
		}
	}

	// Добавляем новую запись для текущего объекта
	lastModified := time.Now().Format(time.RFC3339)
	newRecord := []string{objectKey, fmt.Sprintf("%d", size), contentType, lastModified}
	updatedRecords = append(updatedRecords, newRecord)

	// Перезаписываем CSV файл с обновлёнными записями
	return writeCSV(objectCSVPath, updatedRecords)
}

// readCSV читает все записи из CSV файла
func readCSV(filePath string) ([][]string, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("unable to open objects metadata file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("unable to read objects metadata: %v", err)
	}

	return records, nil
}

// writeCSV перезаписывает CSV файл новыми записями
func writeCSV(filePath string, records [][]string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("unable to open objects metadata file for writing: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Запись всех строк в файл
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("unable to write object metadata: %v", err)
		}
	}

	return nil
}
