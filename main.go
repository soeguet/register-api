package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type RequestValues struct {
	Euro200 [5]int `json:"euro200"`
	Euro100 [5]int `json:"euro100"`
	Euro50  [5]int `json:"euro50"`
	Euro20  [5]int `json:"euro20"`
	Euro10  [5]int `json:"euro10"`
	Euro5   [5]int `json:"euro5"`
	Euro2   [5]int `json:"euro2"`
	Euro1   [5]int `json:"euro1"`
	Cent50  [5]int `json:"cent50"`
	Cent20  [5]int `json:"cent20"`
	Cent10  [5]int `json:"cent10"`
	Cent5   [5]int `json:"cent5"`
	Cent2   [5]int `json:"cent2"`
	Cent1   [5]int `json:"cent1"`
}

type RequestValidation struct {
	TargetValue string `json:"targetValue"`
}

type RequestPayload struct {
	RequestValidation RequestValidation `json:"requestValidation"`
	RequestValues     RequestValues     `json:"requestValues"`
	PayloadType       int               `json:"payloadType"`
}

type ResponseValues struct {
	TotalValue      string `json:"totalValue"`
	DifferenceValue string `json:"differenceValue"`
}

type ResponsePayload struct {
	ResponseValues ResponseValues `json:"responseValues"`
	PayloadType    int            `json:"payloadType"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Received request")
	fmt.Println(r.Body)
	var payload RequestPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Println("Error decoding payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received payload: %+v", payload)

	responsePayload := calculateTotalValue(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsePayload)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func sumArray(arr [5]int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func calculateTotalValue(request RequestPayload) ResponsePayload {
	totalValue := float64(sumArray(request.RequestValues.Euro200)*200) +
		float64(sumArray(request.RequestValues.Euro100)*100) +
		float64(sumArray(request.RequestValues.Euro50)*50) +
		float64(sumArray(request.RequestValues.Euro20)*20) +
		float64(sumArray(request.RequestValues.Euro10)*10) +
		float64(sumArray(request.RequestValues.Euro5)*5) +
		float64(sumArray(request.RequestValues.Euro2)*2) +
		float64(sumArray(request.RequestValues.Euro1)*1) +
		float64(sumArray(request.RequestValues.Cent50))*0.5 +
		float64(sumArray(request.RequestValues.Cent20))*0.2 +
		float64(sumArray(request.RequestValues.Cent10))*0.1 +
		float64(sumArray(request.RequestValues.Cent5))*0.05 +
		float64(sumArray(request.RequestValues.Cent2))*0.02 +
		float64(sumArray(request.RequestValues.Cent1))*0.01

	targetValueAsFloat, _ := strconv.ParseFloat(request.RequestValidation.TargetValue, 64)

	differenceValue := totalValue - targetValueAsFloat

	differenceValueAsStr := strconv.FormatFloat(differenceValue, 'f', 2, 64)

	valueAsStr := strconv.FormatFloat(totalValue, 'f', 2, 64)

	return ResponsePayload{PayloadType: 2, ResponseValues: ResponseValues{TotalValue: valueAsStr, DifferenceValue: differenceValueAsStr}}
}

func main() {
	http.HandleFunc("/api/v1/calculate", corsMiddleware(handleRequest))
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
