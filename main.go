package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type BoxValues struct {
	Euro2  [1]int `json:"euro2"`
	Euro1  [1]int `json:"euro1"`
	Cent50 [1]int `json:"cent50"`
	Cent20 [1]int `json:"cent20"`
	Cent10 [1]int `json:"cent10"`
	Cent5  [1]int `json:"cent5"`
	Cent2  [1]int `json:"cent2"`
	Cent1  [1]int `json:"cent1"`
}

type RollValues struct {
	Euro2  [2]int `json:"euro2"`
	Euro1  [2]int `json:"euro1"`
	Cent50 [2]int `json:"cent50"`
	Cent20 [2]int `json:"cent20"`
	Cent10 [2]int `json:"cent10"`
	Cent5  [2]int `json:"cent5"`
	Cent2  [2]int `json:"cent2"`
	Cent1  [2]int `json:"cent1"`
}

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
	BoxValues         BoxValues         `json:"boxValues"`
	RollValues        RollValues        `json:"rollValues"`
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

func Add(x, y int) int {
	return x + y
}

func handlePayload(w http.ResponseWriter, r *http.Request) (RequestPayload, error) {
	var payload RequestPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Println("Error decoding payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return payload, err
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Only POST method is accepted")
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	payload, err := handlePayload(w, r)
	if err != nil {
		return
	}

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

func SumArray(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func calculateDailyValues(dailyValues RequestValues) float64 {
	return float64(SumArray(dailyValues.Euro200[:])*200) +
		float64(SumArray(dailyValues.Euro100[:])*100) +
		float64(SumArray(dailyValues.Euro50[:])*50) +
		float64(SumArray(dailyValues.Euro20[:])*20) +
		float64(SumArray(dailyValues.Euro10[:])*10) +
		float64(SumArray(dailyValues.Euro5[:])*5) +
		float64(SumArray(dailyValues.Euro2[:])*2) +
		float64(SumArray(dailyValues.Euro1[:])*1) +
		float64(SumArray(dailyValues.Cent50[:]))*0.5 +
		float64(SumArray(dailyValues.Cent20[:]))*0.2 +
		float64(SumArray(dailyValues.Cent10[:]))*0.1 +
		float64(SumArray(dailyValues.Cent5[:]))*0.05 +
		float64(SumArray(dailyValues.Cent2[:]))*0.02 +
		float64(SumArray(dailyValues.Cent1[:]))*0.01
}

func calculateRollValues(rollValues RollValues) float64 {
	return float64(SumArray(rollValues.Euro2[:])*2*25) +
		float64(SumArray(rollValues.Euro1[:])*1*25) +
		float64(SumArray(rollValues.Cent50[:]))*0.5*40 +
		float64(SumArray(rollValues.Cent20[:]))*0.2*40 +
		float64(SumArray(rollValues.Cent10[:]))*0.1*40 +
		float64(SumArray(rollValues.Cent5[:]))*0.05*50 +
		float64(SumArray(rollValues.Cent2[:]))*0.02*50 +
		float64(SumArray(rollValues.Cent1[:]))*0.01*50
}

func calculateBoxValues(boxValues BoxValues) float64 {
	return float64(SumArray(boxValues.Euro2[:])*2*3*25) +
		float64(SumArray(boxValues.Euro1[:])*1*3*25) +
		float64(SumArray(boxValues.Cent50[:]))*0.5*3*40 +
		float64(SumArray(boxValues.Cent20[:]))*0.2*3*40 +
		float64(SumArray(boxValues.Cent10[:]))*0.1*3*40 +
		float64(SumArray(boxValues.Cent5[:]))*0.05*3*40 +
		float64(SumArray(boxValues.Cent2[:]))*0.02*5*50 +
		float64(SumArray(boxValues.Cent1[:]))*0.01*5*50
}

func calculateTotalValue(request RequestPayload) ResponsePayload {
	// calculate intermediate values
	totalValue := calculateDailyValues(request.RequestValues)
	boxValues := calculateBoxValues(request.BoxValues)
	rollValues := calculateRollValues(request.RollValues)
	targetValueAsFloat, _ := strconv.ParseFloat(request.RequestValidation.TargetValue, 64)

	// calculate diff value
	differenceValue := totalValue + boxValues + rollValues - targetValueAsFloat

	// convert to strings
	differenceValueAsStr := strconv.FormatFloat(differenceValue, 'f', 2, 64)
	valueAsStr := strconv.FormatFloat(totalValue+boxValues+rollValues, 'f', 2, 64)

	// response
	return ResponsePayload{
		PayloadType: 2,
		ResponseValues: ResponseValues{
			TotalValue:      valueAsStr,
			DifferenceValue: differenceValueAsStr,
		},
	}
}

func main() {
	http.HandleFunc("/api/v1/calculate", corsMiddleware(handleRequest))
	log.Println("Server starting on port 8002...")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
