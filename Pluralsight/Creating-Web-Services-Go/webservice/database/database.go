package database

import (
    "time"
    "database/sql"
    "log"
)

var DbConn *sql.DB

func SetupDatabase() {
    var err error
    DbConn, err = sql.Open("mysql", "root:sqladmin@tcp(127.0.0.1:3306)/inventorydb")
    DbConn.SetMaxOpenConns(4)
    DbConn.SetMaxIdleConns(4)
    DbConn.SetConnMaxLifetime(60 * time.Second)
    if err != nil {
        log.Fatal(err)
    }
}

