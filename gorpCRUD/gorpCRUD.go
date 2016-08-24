package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

type Post struct {
	// `db` tag lets you specify the column name if it differs from the struct field.
	Id      int64 `db:"post_id"`
	Created time.Time
	Title   string
	Body    string `db:"article_body"`
}

func newPost(title, body string) Post {
	return Post{
		Created: time.Now().UTC(),
		Title:   title,
		Body:    body,
	}
}

func main() {
	// Initialize the DbMap.
	dbmap := initDb()
	defer dbmap.Db.Close()

	// Delete any existing rows.
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed.")

	// Create two posts.
	p1 := newPost("Go 1.1 released!", "Lorem ipsum lorem ipsum")
	p2 := newPost("Go 1.2 released!", "Lorem ipsum lorem ipsum")

	// Insert rows - auto increment PKs will be set properly after the insert.
	err = dbmap.Insert(&p1, &p2)
	checkErr(err, "Insert failed.")

	// Use convenience SelectInt.
	count, err := dbmap.SelectInt("SELECT COUNT(*) FROM posts")
	checkErr(err, "Count failed.")
	log.Println("Rows after inserting:", count)

	// Update a row.
	p2.Title = "Go 1.2 is better than ever."
	count, err = dbmap.Update(&p2)
	checkErr(err, "Update failed.")
	log.Println("Rows updated:", count)

	// Fetch one row - note use of "post_id" instead of "Id" since column is aliased.
	// Postgres users should use $1 instead of ? placeholders.
	err = dbmap.SelectOne(&p2, "SELECT * FROM posts WHERE post_id=$1", p2.Id)
	checkErr(err, "SelectOne failed.")
	log.Println("p2 row:", p2)

	// Fetch all rows.
	var posts []Post
	_, err = dbmap.Select(&posts, "SELECT * FROM posts ORDER BY post_id")
	checkErr(err, "Select failed.")
	log.Println("All rows:")
	for x, p := range posts {
		log.Printf("    %d: %v\n", x, p)
	}

	// Delete row by PK.
	count, err = dbmap.Delete(&p1)
	checkErr(err, "Delete failed.")
	log.Println("Rows deleted:", count)

	// Delete row manually via Exec.
	_, err = dbmap.Exec("DELETE FROM posts WHERE post_id=$1", p2.Id)
	checkErr(err, "Exec failed.")

	// Confirm count is zero.
	count, err = dbmap.SelectInt("SELECT COUNT(*) FROM posts")
	checkErr(err, "Count failed.")
	log.Println("Row count - should be zero:", count)

	log.Println("Done!")
}

const (
	DB_USER     = "penguin"
	DB_PASSWORD = "penguin"
	DB_NAME     = "penguin"
)

func initDb() *gorp.DbMap {
	// Connect to DB using standard Go database/sql API.
	// Use whatever database/sql driver you wish.
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err, "sql.Open failed.")

	// Construct a gorp DbMap.
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// Add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK.
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

	// Create the table. In a production system you'd generally
	// use a migration tool, or create the tables via scripts.
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed.")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
