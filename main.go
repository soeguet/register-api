package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type RequestValues struct {
	Euro200 int `json:"euro200"`
	Euro100 int `json:"euro100"`
	Euro50  int `json:"euro50"`
	Euro20  int `json:"euro20"`
	Euro10  int `json:"euro10"`
	Euro5   int `json:"euro5"`
	Euro2   int `json:"euro2"`
	Euro1   int `json:"euro1"`
	Cent50  int `json:"cent50"`
	Cent20  int `json:"cent20"`
	Cent10  int `json:"cent10"`
	Cent5   int `json:"cent5"`
	Cent2   int `json:"cent2"`
	Cent1   int `json:"cent1"`
}

type RequestPayload struct {
	PayloadType   int           `json:"payloadType"`
	RequestValues RequestValues `json:"requestValues"`
}

type ResponseValues struct {
	TotalValue string `json:"totalValue"`
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

	var payload RequestPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received payload: %+v", payload)

	responsePayload := calculateTotalValue(payload.RequestValues)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsePayload)
}

func calculateTotalValue(request RequestValues) ResponsePayload {
	totalValue := float64(request.Euro200*200) +
		float64(request.Euro100*100) +
		float64(request.Euro50*50) +
		float64(request.Euro20*20) +
		float64(request.Euro10*10) +
		float64(request.Euro5*5) +
		float64(request.Euro2*2) +
		float64(request.Euro1*1) +
		float64(request.Cent50)*0.5 +
		float64(request.Cent20)*0.2 +
		float64(request.Cent10)*0.1 +
		float64(request.Cent5)*0.05 +
		float64(request.Cent2)*0.02 +
		float64(request.Cent1)*0.01

	valueAsStr := strconv.FormatFloat(totalValue, 'f', 2, 64)

	return ResponsePayload{PayloadType: 2, ResponseValues: ResponseValues{TotalValue: valueAsStr}}
}

func main() {
	http.HandleFunc("/calculate", handleRequest)
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
