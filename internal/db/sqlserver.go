package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

var Conn *sql.DB

func Connect() error {
	// SQL Server Docker connection string
	//connString := "sqlserver://sa:my_view_898@127.0.0.1:1433?database=unibazar&encrypt=disable"
	connString := "sqlserver://@127.0.0.1:1477?database=unibazar&trusted_connection=yes"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return err
	}

	// Test connection
	if err := db.PingContext(context.Background()); err != nil {
		return err
	}

	Conn = db
	fmt.Println("✅ Connected to SQL Server")
	return nil
}
