package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"math/rand"
)

type Parameters struct {
	Featured   string `form:"featured"`
	CategoryID string `form:"category_id"`
	BrandID    string `form:"brand_id"`
	StartPrice *int   `form:"start_price"`
	EndPrice   *int   `form:"end_price"`
	Status     string `form:"status"`
	Month      string `form:"month"`
	Key        string `form:"key"`
}

func ProductQueryParameterToMap(P Parameters) string {
	querystring := ""

	if P.Month != "" {
		if querystring != "" {
			querystring = querystring + " AND EXTRACT(MONTH from products.created_at) = " + P.Month

		} else {
			querystring = "EXTRACT(MONTH from products.created_at) = " + P.Month
		}
	}

	if P.CategoryID != "" {
		if querystring != "" {
			querystring = querystring + " AND category_id IN (" + P.CategoryID + ")"

		} else {
			querystring = "category_id IN (" + P.CategoryID + ")"
		}
	}
	if P.BrandID != "" {
		if querystring != "" {
			querystring = querystring + " AND brand_id IN (" + P.BrandID + ")"

		} else {
			querystring = "brand_id IN (" + P.BrandID + ")"
		}
	}
	if P.Featured != "" {
		if querystring != "" {
			querystring = querystring + " AND products.featured = " + P.Featured

		} else {
			querystring = "products.featured = " + P.Featured
		}
	}
	if P.Status != "" {
		if querystring != "" {
			querystring = querystring + " AND products.status = '" + P.Status + "'"

		} else {
			querystring = "products.status = '" + P.Status + "'"
		}
	}

	if P.StartPrice != nil && P.EndPrice != nil {
		if querystring != "" {
			querystring = querystring + " AND " + fmt.Sprintf("price >= %d AND price <= %d", *P.StartPrice, *P.EndPrice)

		} else {
			querystring = fmt.Sprintf("price >= %d AND price <= %d", *P.StartPrice, *P.EndPrice)
		}

	} else if P.StartPrice != nil && P.EndPrice == nil {
		if querystring != "" {
			querystring = querystring + " AND " + fmt.Sprintf("price >= %d", *P.StartPrice)

		} else {
			querystring = fmt.Sprintf("price >= %d", *P.StartPrice)
		}

	} else if P.StartPrice == nil && P.EndPrice != nil {
		if querystring != "" {
			querystring = querystring + " AND " + fmt.Sprintf("price <= %d", *P.EndPrice)

		} else {
			querystring = fmt.Sprintf("price <= %d", *P.EndPrice)
		}
	}
	return querystring
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

// Decode Base64 string to []byte
func DecodeBase64Image(base64String string) ([]byte, error) {
	decodedImage, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}
	return decodedImage, nil
}

func EncodeImageToBase64(imageBytes []byte) string {
	encodedString := base64.StdEncoding.EncodeToString(imageBytes)
	return encodedString
}
