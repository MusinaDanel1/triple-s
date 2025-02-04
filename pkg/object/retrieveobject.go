package object

import (
	"net/http"
	"os"
	"path/filepath"
)

// RetrieveObjectHandler обрабатывает запрос на получение объекта из бакета.
func RetrieveObjectHandler(w http.ResponseWriter, r *http.Request, bucketDir, bucketName, objectKey string) {
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

	// 4. Определяем Content-Type на основе расширения файла
	contentType := http.DetectContentType([]byte{})
	ext := filepath.Ext(objectKey)
	if ext == ".png" {
		contentType = "image/png"
	} else if ext == ".jpg" || ext == ".jpeg" {
		contentType = "image/jpeg"
	} else if ext == ".txt" {
		contentType = "text/plain"
	} else {
		contentType = "application/octet-stream" // по умолчанию для неизвестных типов
	}

	// 5. Чтение содержимого объекта
	file, err := os.Open(objectPath)
	if err != nil {
		http.Error(w, "500 Internal Server Error: Unable to open object", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "500 Internal Server Error: Unable to get object info", http.StatusInternalServerError)
		return
	}

	// Читаем содержимое файла в память
	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		http.Error(w, "500 Internal Server Error: Unable to read object", http.StatusInternalServerError)
		return
	}

	// 6. Устанавливаем заголовки и возвращаем данные объекта
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
