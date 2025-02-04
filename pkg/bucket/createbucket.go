package bucket

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// Bucket представляет структуру ведра
type Bucket struct {
	XMLName          xml.Name `xml:"Bucket"`
	Name             string   `xml:"Name"`
	CreationTime     string   `xml:"CreationTime"`
	LastModifiedTime string   `xml:"LastModifiedTime"`
	Status           string   `xml:"Status"`
}

// validateBucketName проверяет имя ведра на соответствие правилам
func validateBucketName(bucketName string) (bool, string) {
	const namePattern = `^[a-z0-9.-]{3,63}$`

	// Проверка шаблона имени
	if matched, _ := regexp.MatchString(namePattern, bucketName); !matched {
		return false, "bucket name must be between 3 and 63 characters and can only contain lowercase letters, numbers, hyphens, and periods"
	}

	// Дополнительные правила проверки
	if bucketName[0] == '-' || bucketName[len(bucketName)-1] == '-' {
		return false, "bucket name must not begin or end with a hyphen"
	}
	if regexp.MustCompile(`--|\.\.`).MatchString(bucketName) {
		return false, "bucket name must not contain two consecutive periods or dashes"
	}
	if net.ParseIP(bucketName) != nil {
		return false, "bucket name must not be formatted as an IP address"
	}

	return true, ""
}

// isBucketNameUnique проверяет уникальность имени ведра по данным в CSV
func isBucketNameUnique(bucketName, csvFilePath string) (bool, error) {
	file, err := os.Open(csvFilePath)
	if os.IsNotExist(err) {
		// Если файл не существует, то имя уникально
		return true, nil
	} else if err != nil {
		return false, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	// Проверка на существование ведра с таким именем
	for _, record := range records {
		if record[0] == bucketName {
			return false, nil
		}
	}

	return true, nil
}

// createBucket создает ведро, директорию и записывает информацию о ведре в CSV
func createBucket(bucketName, csvFilePath, dataDir string) (Bucket, error) {
	// 1. Проверка имени ведра
	valid, msg := validateBucketName(bucketName)
	if !valid {
		return Bucket{}, fmt.Errorf("400 bad request: %s", msg)
	}

	// 2. Проверка уникальности имени ведра
	unique, err := isBucketNameUnique(bucketName, csvFilePath)
	if err != nil {
		return Bucket{}, fmt.Errorf("error checking bucket uniqueness: %v", err)
	}
	if !unique {
		return Bucket{}, fmt.Errorf("409 conflict: bucket name already exists")
	}

	// 3. Создание подкаталога для ведра
	bucketPath := filepath.Join(dataDir, bucketName)
	err = os.Mkdir(bucketPath, 0o755)
	if err != nil {
		return Bucket{}, fmt.Errorf("error creating bucket directory: %v", err)
	}

	// 5. Создаем информацию о ведре
	creationTime := time.Now().Format(time.RFC3339)
	lastModifiedTime := creationTime
	status := "active"

	bucket := Bucket{
		Name:             bucketName,
		CreationTime:     creationTime,
		LastModifiedTime: lastModifiedTime,
		Status:           status,
	}

	// 6. Запись информации о ведре в CSV файл
	file, err := os.OpenFile(csvFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return Bucket{}, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	bucketMetadata := []string{bucket.Name, bucket.CreationTime, bucket.LastModifiedTime, bucket.Status}
	if err := writer.Write(bucketMetadata); err != nil {
		return Bucket{}, fmt.Errorf("error writing to CSV file: %v", err)
	}

	// Возвращаем созданную корзину
	return bucket, nil
}

// CreateBucketHandler обрабатывает HTTP-запросы на создание ведра
func CreateBucketHandler(w http.ResponseWriter, r *http.Request, dataDir, bucketName string) {
	if bucketName == "" {
		http.Error(w, "400 bad request: missing bucket name in the URL", http.StatusBadRequest)
		return
	}

	// Путь к CSV файлу и директории для хранения данных
	csvFilePath := filepath.Join(dataDir, "buckets.csv")

	// Вызов функции createBucket для создания ведра
	bucket, err := createBucket(bucketName, csvFilePath, dataDir)
	if err != nil {
		if err.Error()[:3] == "400" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error()[:3] == "409" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "500 internal server error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Устанавливаем заголовок Content-Type для XML
	w.Header().Set("Content-Type", "application/xml")

	// Кодируем структуру ведра в XML
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(bucket); err != nil {
		http.Error(w, "500 internal server error: unable to encode XML", http.StatusInternalServerError)
	}
}
