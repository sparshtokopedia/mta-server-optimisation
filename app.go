package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type HostNameIPStatus struct {
	IP     string
	Status bool
}

type GetHostNameResponse struct {
	Result []string `json:"result"`
	Status string   `json:"status"`
	Error  error    `json:"error"`
}

var ipMap map[string][]HostNameIPStatus

func main() {
	err := loadMockData()
	if err != nil {
		log.Fatal("Failed to load mock data: ", err)
	}

	http.HandleFunc("/mta-hosting-optimizer", getInstanceName)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getInstanceName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	thresholdEnv := GoDotEnvVariable("X")
	// thresholdEnv := getEnv("X","1")
	threshold, err := strconv.Atoi(thresholdEnv)
	if err != nil {
		log.Println("Error converting string to int")
		json.NewEncoder(w).Encode(GetHostNameResponse{
			Result: nil,
			Status: "Error",
			Error:  err,
		})
		return
	}

	result := getInefficientInstance(threshold)

	json.NewEncoder(w).Encode(result)
}

func getInefficientInstance(threshold int) []string {
	InefficientInstance := make([]string, 0)

	for key, val := range ipMap {
		count := 0
		for _, ipStatus := range val {
			if ipStatus.Status {
				count++
			}
		}
		if count <= threshold {
			InefficientInstance = append(InefficientInstance, key)
		}
	}

	return InefficientInstance
}

func loadMockData() error {
	ips := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"}
	hostNames := []string{"mta-prod-1", "mta-prod-1", "mta-prod-2", "mta-prod-2", "mta-prod-2", "mta-prod-3"}
	actives := []bool{true, false, true, true, false, false}

	ipMap = make(map[string][]HostNameIPStatus)

	for idx := 0; idx < len(ips); idx++ {
		ipMap[hostNames[idx]] = append(ipMap[hostNames[idx]], HostNameIPStatus{
			IP:     ips[idx],
			Status: actives[idx],
		})
	}

	return nil
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
