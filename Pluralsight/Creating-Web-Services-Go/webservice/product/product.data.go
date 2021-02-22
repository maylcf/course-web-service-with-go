package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products ...")

	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d products loaded... \n", len(productMap.m))
}

func loadProductMap() (map[int]Product, err) {
	fileName := "products.json"

	// Verify if file exist
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	// Read data in the file
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err := json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}

	// Convert products list to a map
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}

	return prodMap, nil
}

func getProduct(productID int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()

	if product, ok := productMap.m[productID]; ok {
		return &product
	}

	return nil
}

func removeProduct(productID int) {
	productMap.RLock()
	defer productMap.RUnlock()
	delete(productMap.m, productID)
}

func getProductList() []Product {
	productMap.RLock()

	products := make([]Product, 0, len(productMap.m))

	for _, value := range productMap.m {
		products = append(products, value)
	}

	productMap.Unlock()

	return products
}

func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}

	for key := range productMap.m {
		productIds = append(productIds, key)
	}

	productMap.RUnlock()
	sort.Ints(productIds)

	return productIds
}

func getNextProductId() int {
	productIDs := getProductIds()
	return productIDs[len(productIDs)-1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	// if the product id is set, update, otherwise add
	addOrUpdateId := -1

	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)

		// if it exists, replace it, otherwise return error
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%id] doesn't exist", product.ProductID)
		}

		addOrUpdateId = product.ProductID
	} else {
		addOrUpdateId = getNextProductId()
		product.ProductID = addOrUpdateId
	}

	productMap.Lock()
	productMap.m[addOrUpdateId] = product
	productMap.Unlock()

	return addOrUpdateId, nil
}
