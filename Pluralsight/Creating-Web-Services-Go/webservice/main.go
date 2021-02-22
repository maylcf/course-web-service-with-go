package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/maylcf/learning-go/product"
)

type fooHandler struct {
	Message string
}

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func init() {
	productsJSON := `[{
		"productId": 1,
		"manufacturer": "Johns-Jenkins",
		"sku": "abs123",
		"upc": "12345",
		"pricePerUnit": "430.50",
		"quantityOnHand": 9703,
		"productName": "sticky note"
	},
	{
		"productId": 2,
		"manufacturer": "Hessel, Schimmel and Feeney",
		"sku": "abs456",
		"upc": "67890",
		"pricePerUnit": "280.50",
		"quantityOnHand": 9710,
		"productName": "leg warmers"
	},
	{
		"productId": 3,
		"manufacturer": "Swaniawski",
		"sku": "abc789",
		"upc": "90543",
		"pricePerUnit": "100.00",
		"quantityOnHand": 5905,
		"productName": "lamp shade"
	}
]`

	err := json.Unmarshal([]byte(productsJSON), &productList)

	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1

	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}

	return highestID + 1
}

func findProductById(productID int) (*Product, int) {
	for i, product := range productList {
		if product.ProductID == productID {
			return &product, i
		}
	}
	return nil, 0
}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before handler; middleware start")
		start := time.Now()

		handler.ServeHTTP(w, r)
		fmt.Printf("middleware finished; %s", time.Since(start))
	})
}

const apiBasePath = "/api"

func main() {

	product.SetupRoutes(apiBasePath)
	productListHandler := http.HandlerFunc(productsHandler)
	productItemHandler := http.HandlerFunc(productHandler)

	http.Handle("/products", middlewareHandler(productListHandler))
	http.Handle("/products/", middlewareHandler(productItemHandler))

	http.ListenAndServe(":5000", nil)
}
