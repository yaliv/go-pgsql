package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

type Invoice struct {
	Id       int64
	Created  int64
	Updated  int64
	Memo     string
	PersonId int64
}

type Person struct {
	Id      int64
	Created int64
	Updated int64
	FName   string
	LName   string
}

// Define a type for your join
// It *must* contain all the columns in your SELECT statement
//
// The names here should match the aliased column names you specify
// in your SQL - no additional binding work required. Simple!
//
type InvoicePersonView struct {
	InvoiceId int64
	PersonId  int64
	Memo      string
	FName     string
}

func main() {
	// Initialize the DbMap.
	dbmap := initDb()
	defer dbmap.Db.Close()

	// Find a person.
	var p1 *Person
	err := dbmap.SelectOne(&p1, "SELECT * FROM persons WHERE id=1")

	// Create it if not found.
	if err != nil {
		p1 = &Person{0, 0, 0, "Jerry", "Gray"}
		err := dbmap.Insert(p1)
		checkErr(err, "Insert person failed.")
	}

	// Notice how we can connect p1.Id to the invoice easily.
	inv1 := &Invoice{0, 0, 0, "dropship order", p1.Id}
	err = dbmap.Insert(inv1)
	checkErr(err, "Insert invoice failed.")

	// Run your query.
	query := "SELECT i.Id InvoiceId, p.Id PersonId, i.Memo, p.FName " +
		"FROM invoices i, persons p " +
		"WHERE i.PersonId = p.Id"

	// Pass a slice to Select().
	var list []InvoicePersonView
	_, err = dbmap.Select(&list, query)
	checkErr(err, "Join failed.")

	lastFromJoin := list[len(list)-1]
	fmt.Println("Last from join:", lastFromJoin)

	// This should test true.
	expected := InvoicePersonView{inv1.Id, p1.Id, inv1.Memo, p1.FName}
	fmt.Println("Expected:", expected)

	if reflect.DeepEqual(lastFromJoin, expected) {
		fmt.Println("Yaayy! My join worked!")
	}
}

const (
	DB_USER     = "penguin"
	DB_PASSWORD = "penguin"
	DB_NAME     = "penguin"
)

func initDb() *gorp.DbMap {
	// Connect to DB using standard Go database/sql API.
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err, "sql.Open failed")

	// Construct a gorp DbMap.
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// Add tables, setting the table names and specifying that
	// their Id property are auto incrementing PK.
	dbmap.AddTableWithName(Invoice{}, "invoices").SetKeys(true, "Id")
	dbmap.AddTableWithName(Person{}, "persons").SetKeys(true, "Id")

	// Create the tables. In a production system you'd generally
	// use a migration tool, or create the tables via scripts.
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
