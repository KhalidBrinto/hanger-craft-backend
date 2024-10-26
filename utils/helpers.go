package utils

import (
	"time"

	"math/rand"
)

func GenerateOrderID() string {
	const charset = "0123456789"
	length := 6
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random alphanumeric string of length 6
	orderID := make([]byte, length)
	for i := range orderID {
		orderID[i] = charset[seededRand.Intn(len(charset))]
	}

	// Return the ID with 'HC' prefix
	return "HC" + string(orderID)
}

func GenerateTransactionID() string {
	const charset = "0123456789"
	length := 8
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random alphanumeric string of length 6
	txID := make([]byte, length)
	for i := range txID {
		txID[i] = charset[seededRand.Intn(len(charset))]
	}

	// Return the ID with 'HC' prefix
	return "INV" + string(txID)
}
