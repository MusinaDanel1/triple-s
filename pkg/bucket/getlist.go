package bucket

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Структура для формирования XML-ответа
type ListAllBucketsResponse struct {
	XMLName xml.Name `xml:"ListAllBucketsResponse"`
	Buckets []Bucket `xml:"Bucket"`
}

// Функция для получения всех ведер из CSV файла
func getAllBucketsFromCSV(csvFilePath string) ([]Bucket, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %v", err)
	}

	var buckets []Bucket
	for _, record := range records {
		bucket := Bucket{
			Name:             record[0],
			CreationTime:     record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

// ListAllBucketsHandler обрабатывает HTTP-запросы на получение списка ведер.
func ListAllBucketsHandler(w http.ResponseWriter, r *http.Request, dataDir string) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	csvFilePath := filepath.Join(dataDir, "buckets.csv")

	buckets, err := getAllBucketsFromCSV(csvFilePath)
	if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем ответ в формате XML
	response := ListAllBucketsResponse{Buckets: buckets}

	// Устанавливаем Content-Type для XML
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)

	// Кодируем ответ в XML и отправляем его клиенту
	err = xml.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "500 Internal Server Error: Failed to encode XML response", http.StatusInternalServerError)
	}
}
