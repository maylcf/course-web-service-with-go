package main

import (
	"net/http"

	"github.com/maylcf/learning-go/product"
	"github.com/maylcf/learning-go/database"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
    database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
