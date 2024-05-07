package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	RollsPerBoxesFive  = 5
	RollsPerBoxesThree = 3

	CoinsPerRollEuro      = 25
	CoinsPerRollBigCent   = 40
	CoinsPerRollSmallCent = 50

	Euro2Value  = 2.0
	Euro1Value  = 1.0
	Cent50Value = 0.5
	Cent20Value = 0.2
	Cent10Value = 0.1
	Cent5Value  = 0.05
	Cent2Value  = 0.02
	Cent1Value  = 0.01
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

// FormatNumber takes a float64 value `value` and formats it as a string with two decimal places.
// The function splits the string representation of `value` into its integer and decimal parts.
// It then checks if the integer part is negative and temporarily removes the negative sign if present.
// The function adds thousand separators to the integer part by iterating over the characters in reverse order
// and inserting a dot separator every three digits.
// The resulting integer part is reversed to correct the order of the digits.
// If the original value was negative, the negative sign is prepended back to the formatted string.
// Finally, the function returns the formatted string by combining the integer part, decimal part, and
// a comma separator between them.
func FormatNumber(value float64) string {
	// convert to a string with two decimal places
	str := fmt.Sprintf("%.2f", value)
	parts := strings.Split(str, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	// check for a negative sign and remove it temporarily
	negative := false
	if integerPart[0] == '-' {
		negative = true
		integerPart = integerPart[1:]
	}

	// adding thousand separators
	n := len(integerPart)
	var withSeparator strings.Builder
	for i := n - 1; i >= 0; i-- {
		withSeparator.WriteByte(integerPart[i])
		if (n-i)%3 == 0 && i != 0 {
			withSeparator.WriteByte('.')
		}
	}

	// reverse to correct the order
	result := reverseString(withSeparator.String())

	// prepend negative sign if the number was negative
	if negative {
		result = "-" + result
	}

	return fmt.Sprintf("%s,%s", result, decimalPart)
}

// reverseString takes a string `s` and reverses its order by swapping characters from the beginning and end of the string.
// It returns the reversed string.
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// HandlePayload decodes the request payload from the HTTP request body into a `RequestPayload` struct.
// If there is an error decoding the payload, it returns an HTTP error response with a 400 status code.
// The error is also wrapped and returned as an error value.
// If the decoding is successful, it returns the decoded payload and a nil error.
func HandlePayload(w http.ResponseWriter, r *http.Request) (RequestPayload, error) {
	var payload RequestPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return payload, fmt.Errorf("error decoding payload: %w", err)
	}
	return payload, nil
}

// handlePOSTRequest handles HTTP POST requests.
// It checks if the request method is POST and returns an error if it's not.
// It then calls the HandlePayload function to decode the request payload.
// If there is an error decoding the payload, handlePOSTRequest returns early.
// It then calls the calculateTotalValue function to calculate the total value based on the payload.
// Finally, it calls the respondWithJSON function to send the response payload as a JSON response.
func handlePOSTRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Only POST method is accepted")
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}
	payload, err := HandlePayload(w, r)
	if err != nil {
		return
	}
	responsePayload := calculateTotalValue(payload)
	respondWithJSON(w, responsePayload)
}

// respondWithJSON sets the "Content-Type" header of the HTTP response to "application/json".
// It also writes the response payload as JSON to the response writer.
// If there is an error encoding the response payload, it returns early without writing anything.
func respondWithJSON(w http.ResponseWriter, responsePayload ResponsePayload) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(responsePayload)
	if err != nil {
		return
	}
}

// corsMiddleware is a middleware function that adds the necessary CORS headers to the HTTP response.
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

// SumArray calculates the sum of all elements in the given integer array.
// It iterates through each element in the array and adds it to the running sum.
// The final sum is returned as an integer value.
func SumArray(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

// CalculateDailyValues calculates the total value of the daily values in euros and cents
// based on the given RequestValues struct.
// The function uses the SumArray helper function to calculate the sum of each array and multiplies
// it by the corresponding euro or cent value.
// The result is the sum of all the calculated values.
func CalculateDailyValues(dailyValues RequestValues) float64 {
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

// CalculateRollValues calculates the total value of the given RollValues struct. It
// multiplies the sum of each array by the corresponding coin value and the number of coins
// per roll. The calculated values are then summed up and returned as a float64 value.
func CalculateRollValues(rollValues RollValues) float64 {
	coinSumEuro2 := float64(SumArray(rollValues.Euro2[:])) * Euro2Value * CoinsPerRollEuro
	coinSumEuro1 := float64(SumArray(rollValues.Euro1[:])) * Euro1Value * CoinsPerRollEuro
	coinSumCent50 := float64(SumArray(rollValues.Cent50[:])) * Cent50Value * CoinsPerRollBigCent
	coinSumCent20 := float64(SumArray(rollValues.Cent20[:])) * Cent20Value * CoinsPerRollBigCent
	coinSumCent10 := float64(SumArray(rollValues.Cent10[:])) * Cent10Value * CoinsPerRollBigCent
	coinSumCent5 := float64(SumArray(rollValues.Cent5[:])) * Cent5Value * CoinsPerRollSmallCent
	coinSumCent2 := float64(SumArray(rollValues.Cent2[:])) * Cent2Value * CoinsPerRollSmallCent
	coinSumCent1 := float64(SumArray(rollValues.Cent1[:])) * Cent1Value * CoinsPerRollSmallCent

	return coinSumEuro2 + coinSumEuro1 + coinSumCent50 + coinSumCent20 + coinSumCent10 + coinSumCent5 + coinSumCent2 + coinSumCent1
}

// calculateIndividualBoxValue calculates the value of an individual box based on the provided array of integers,
// a value multiplier, a times multiplier, and a multiplier factor.
// It first sums up the elements in the array using the SumArray function.
// Then, it multiplies the sum by the value, times, and multiplier constants.
// The final result is returned as a float64 value.
func calculateIndividualBoxValue(arr []int, value, times, multiplier float64) float64 {
	return float64(SumArray(arr)) * value * times * multiplier
}

func CalculateBoxValues(box BoxValues) float64 {
	return calculateIndividualBoxValue(box.Euro2[:], Euro2Value, RollsPerBoxesThree, CoinsPerRollEuro) +
		calculateIndividualBoxValue(box.Euro1[:], Euro1Value, RollsPerBoxesThree, CoinsPerRollEuro) +
		calculateIndividualBoxValue(box.Cent50[:], Cent50Value, RollsPerBoxesThree, CoinsPerRollBigCent) +
		calculateIndividualBoxValue(box.Cent20[:], Cent20Value, RollsPerBoxesThree, CoinsPerRollBigCent) +
		calculateIndividualBoxValue(box.Cent10[:], Cent10Value, RollsPerBoxesThree, CoinsPerRollBigCent) +
		calculateIndividualBoxValue(box.Cent5[:], Cent5Value, RollsPerBoxesThree, CoinsPerRollSmallCent) +
		calculateIndividualBoxValue(box.Cent2[:], Cent2Value, RollsPerBoxesFive, CoinsPerRollSmallCent) +
		calculateIndividualBoxValue(box.Cent1[:], Cent1Value, RollsPerBoxesFive, CoinsPerRollSmallCent)
}

// calculateTotalValue calculates the total value based on the given RequestPayload struct.
// It calls the CalculateValuesForCashCounts function to calculate the intermediate values.
// It converts the differenceValue and totalValue+boxValues+rollValues to strings using strconv.FormatFloat.
// It constructs and returns a ResponsePayload struct with the calculated values.
func calculateTotalValue(request RequestPayload) ResponsePayload {
	totalValue, boxValues, rollValues, differenceValue := CalculateValuesForCashCounts(request)

	// convert to strings
	// differenceValueAsStr := strconv.FormatFloat(differenceValue, 'f', 2, 64)
	differenceValueAsStr := FormatNumber(differenceValue)
	// valueAsStr := strconv.FormatFloat(totalValue+boxValues+rollValues, 'f', 2, 64)
	valueAsStr := FormatNumber(totalValue + boxValues + rollValues)

	// response
	return ResponsePayload{
		PayloadType: 2,
		ResponseValues: ResponseValues{
			TotalValue:      valueAsStr,
			DifferenceValue: differenceValueAsStr,
		},
	}
}

// CalculateValuesForCashCounts calculates the total value, box value, roll value, and difference value
// based on the given RequestPayload struct. It uses the CalculateDailyValues, CalculateBoxValues,
// and CalculateRollValues functions to calculate the intermediate values. It converts the target value
// from string to float64 using strconv.ParseFloat. The difference value is calculated as the difference
// between the sum of total value, box value, and roll value, and the target value as a float64.
// The function returns the calculated total value, box value, roll value, and difference value as float64.
func CalculateValuesForCashCounts(request RequestPayload) (float64, float64, float64, float64) {
	// calculate intermediate values
	totalValue := CalculateDailyValues(request.RequestValues)
	boxValues := CalculateBoxValues(request.BoxValues)
	rollValues := CalculateRollValues(request.RollValues)
	targetValueAsFloat, _ := strconv.ParseFloat(request.RequestValidation.TargetValue, 64)

	// calculate diff value
	differenceValue := totalValue + boxValues + rollValues - targetValueAsFloat
	return totalValue, boxValues, rollValues, differenceValue
}

// main starts the HTTP server and registers the handler functions.
// It listens for requests on the "/api/v1/calculate" endpoint.
// It uses the corsMiddleware function to add the necessary CORS headers.
// If there is an error starting the server, it logs the error and exits.
func main() {
	http.HandleFunc("/api/v1/calculate", corsMiddleware(handlePOSTRequest))
	log.Println("Server starting on port 8002...")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
