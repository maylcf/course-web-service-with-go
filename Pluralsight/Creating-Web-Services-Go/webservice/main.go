package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")

	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	product, listItemIndex := findProductById(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Return a single product
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)

	case http.MethodPost:
		// Update product in list
		var updatedProduct Product

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify product ID
		if updatedProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Update product
		product = &updatedProduct
		productList[listItemIndex] = *product

		// Return success
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)

	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newProduct.ProductID = getNextID()
		productList = append(productList, newProduct)

		w.WriteHeader(http.StatusCreated)
		return
	}
}

// func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(f.Message))
// }

// func barHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("bar called"))
// }

func main() {
	// http.Handle("/foo", &fooHandler{Message: "foo called"})
	// http.HandleFunc("/bar", barHandler)

	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)

	http.ListenAndServe(":5000", nil)
}
