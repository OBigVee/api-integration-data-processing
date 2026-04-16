package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"encoding/json"
	"net/http"
)

type GenderizeResponse struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type SuccessResponse struct {
	Status string `json:"status"`
	Data   struct {
		Name        string  `json:"name"`
		Gender      string  `json:"gender"`
		Probability float64 `json:"probability"`
		SampleSize  int     `json:"sample_size"`
		IsConfident bool    `json:"is_confident"`
		ProcessedAt string  `json:"processed_at"`
	} `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	//baseUrl := "http://localhost" +  addr

	http.HandleFunc("/api/classify", classifyHandler)

	fmt.Printf("Server is live and listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed : %s", err)
	}
}

func classifyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s]  %s", r.Method, r.URL.String())

	// set CORS Header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	const Genderize_URL string = "https://api.genderize.io/?name="

	//validate query params
	name := r.URL.Query().Get("name")
	if name == "" {
		sendError(w, "Missing or empty name parameter", http.StatusBadRequest)
		return
	}

	for _, char := range name {
		if !unicode.IsLetter(char) && char != ' ' && char != '-' {
			sendError(w, "name is not a string", http.StatusUnprocessableEntity)
			return
		}
	}

	// Call Genderize API
	resp, err := httpClient.Get(Genderize_URL + name)
	if err != nil {
		
		sendError(w, "Upstream or server failure", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var gData GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&gData); err != nil {
		sendError(w, "Upstream or server failure", http.StatusInternalServerError)
		return
	}

	if gData.Gender == "" || gData.Count == 0 {
		sendError(w, "No prediction available for the provided name", http.StatusOK)
		return
	}

	// Rule : Probability >= 0.7 and Sample Size >= 100
	isConfident := gData.Probability >= 0.7 && gData.Count >= 100

	// Prepare response
	finalResp := SuccessResponse{
		Status: "success",
	}

	finalResp.Data.Name = gData.Name
	finalResp.Data.Gender = gData.Gender
	finalResp.Data.Probability = gData.Probability
	finalResp.Data.SampleSize = gData.Count
	finalResp.Data.IsConfident = isConfident
	finalResp.Data.ProcessedAt = time.Now().UTC().Format(time.RFC3339)

	// Send response
	json.NewEncoder(w).Encode(finalResp)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("REJECTED: Code %d - %s", statusCode, message)

	w.WriteHeader(statusCode)
	errorResp := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	json.NewEncoder(w).Encode(errorResp)
}
