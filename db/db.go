package db

import (
	"database/sql"
	"fmt"
	. "login_jwt/helpers"

	_ "github.com/go-sql-driver/mysql"
)

func CreateCon() *sql.DB {
	cfg := ReadAppConfig()
	db, err := sql.Open("mysql", cfg["DBUser"]+":"+cfg["DBPass"]+"@tcp("+cfg["DBHost"]+":"+cfg["DBPort"]+")/"+cfg["DBDatabase"]+"?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}
	// defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("db is NOT Connected")
		fmt.Println(err.Error())
	}
	return db
}
