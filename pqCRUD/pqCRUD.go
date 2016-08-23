package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "penguin"
	DB_PASSWORD = "penguin"
	DB_NAME     = "penguin"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	var lastInsertId int

	fmt.Println("# Inserting values")

	err = db.QueryRow("INSERT INTO userinfo (username,departname,created) "+
		"VALUES ($1,$2,$3), ($4,$5,$6), ($7,$8,$9), ($10,$11,$12) RETURNING uid;",

		"skipper", "General Management", "2016-01-25",
		"kowalski", "Production", "2016-03-27",
		"rico", "Warehouse", "2016-06-06",
		"private", "Human Resource", nil,
	).Scan(&lastInsertId)
	checkErr(err)

	fmt.Println("last inserted id =", lastInsertId)

	fmt.Println("# Updating")

	stmt, err := db.Prepare("UPDATE userinfo SET username=$1 WHERE uid=$2")
	checkErr(err)

	res, err := stmt.Exec("skippersir", lastInsertId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")

	fmt.Println("# Querying")

	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	fmt.Println("uid | username | department | created")
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created interface{}
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)

		if created != nil {
			created = created.(time.Time)
		}

		fmt.Printf("%3v | %8v | %10v | %v\n", uid, username, department, created)
	}

	fmt.Println("# Deleting")

	stmt, err = db.Prepare("DELETE FROM userinfo WHERE uid=$1")
	checkErr(err)

	res, err = stmt.Exec(lastInsertId)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
