package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Receipt structure defines the fields of a receipt.
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Item structure defines fields for an individual item on the receipt.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// ReceiptPoints stores the ID and points associated with a receipt.
type ReceiptPoints struct {
	ID     string
	Points int
}

var receipts = make(map[string]ReceiptPoints)

// processReceipt handles the POST /receipts/process endpoint.
func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid receipt format", http.StatusBadRequest)
		return
	}

	// Generate a consistent UUID based on receipt content.
	id := generateUUID(receipt)

	// If the receipt is not already processed, calculate points and store.
	if _, exists := receipts[id]; !exists {
		points := calculatePoints(receipt)
		receipts[id] = ReceiptPoints{ID: id, Points: points}
	}

	// Respond with the receipt's ID.
	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getPoints handles the GET /receipts/{id}/points endpoint.
func getPoints(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if receiptPoints, exists := receipts[id]; exists {
		response := map[string]int{"points": receiptPoints.Points}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Receipt not found", http.StatusNotFound)
	}
}

// calculatePoints calculates points for a given receipt based on rules.
func calculatePoints(receipt Receipt) int {
	points := 0

	// 1 point for every alphanumeric character in retailer name.
	points += countAlphanumeric(receipt.Retailer)

	// 50 points if total is a round dollar amount.
	if isRoundDollar(receipt.Total) {
		points += 50
	}

	// 25 points if total is a multiple of 0.25.
	if isMultipleOf(receipt.Total, 0.25) {
		points += 25
	}

	// 5 points for every two items.
	points += (len(receipt.Items) / 2) * 5

	// Points for item descriptions whose trimmed length is a multiple of 3.
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			itemPrice, _ := parsePrice(item.Price)
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd.
	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day()%2 != 0 {
		points += 6
	}

	// 10 points if the time is between 2:00 PM and 4:00 PM.
	time, _ := time.Parse("15:04", receipt.PurchaseTime)
	if time.Hour() == 14 {
		points += 10
	}

	return points
}

// generateUUID creates a consistent UUID for a receipt based on its content.
func generateUUID(receipt Receipt) string {
	receiptBytes, _ := json.Marshal(receipt)
	hash := sha256.Sum256(receiptBytes)
	return uuid.NewSHA1(uuid.NameSpaceOID, hash[:]).String()
}

// Helper function to count alphanumeric characters in a string.
func countAlphanumeric(s string) int {
	reg := regexp.MustCompile("[a-zA-Z0-9]")
	return len(reg.FindAllString(s, -1))
}

// Helper function to check if a total is a round dollar amount.
func isRoundDollar(total string) bool {
	price, _ := parsePrice(total)
	return math.Mod(price, 1.0) == 0
}

// Helper function to check if a total is a multiple of a given value.
func isMultipleOf(total string, factor float64) bool {
	price, _ := parsePrice(total)
	return math.Mod(price, factor) == 0
}

// Helper function to parse a string price into a float64.
func parsePrice(priceStr string) (float64, error) {
	return strconv.ParseFloat(priceStr, 64)
}

// main sets up the server and routes.
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
