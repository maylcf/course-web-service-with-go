package main

import (
	"net/http"

	"github.com/maylcf/learning-go/product"
	"github.com/maylcf/learning-go/database"
	"github.com/maylcf/learning-go/receipt"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
    database.SetupDatabase()
    receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
