package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123123@tcp(127.0.0.1:3306)/gotest")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	stmt, err := db.Prepare("DROP TABLE todo;")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("TODOS REMOVED")
	}
}
