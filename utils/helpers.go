package utils

import (
	"fmt"
	"time"

	"math/rand"
)

type Parameters struct {
	Featured   *bool   `form:"featured"`
	CategoryID *uint   `form:"category_id"`
	StartPrice *int    `form:"start_price"`
	EndPrice   *int    `form:"end_price"`
	Status     *string `form:"status"`
}

func ProductQueryParameterToMap(P Parameters) (map[string]interface{}, string) {
	querystring := ""
	QueryMap := make(map[string]interface{})
	if P.Featured != nil {
		QueryMap["products.featured"] = P.Featured
	}
	if P.Status != nil {
		QueryMap["products.status"] = P.Featured
	}
	if P.CategoryID != nil {
		QueryMap["products.category_id"] = P.CategoryID
	}

	if P.StartPrice != nil && P.EndPrice != nil {
		querystring = fmt.Sprintf("price >= %d AND price <= %d", *P.StartPrice, *P.EndPrice)
	} else if P.StartPrice != nil && P.EndPrice == nil {
		querystring = fmt.Sprintf("price >= %d", *P.StartPrice)

	} else if P.StartPrice == nil && P.EndPrice != nil {
		querystring = fmt.Sprintf("price <= %d", *P.EndPrice)
	}

	return QueryMap, querystring
}

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
