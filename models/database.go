package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

var dbmap *gorp.DbMap

func init() {
	var err error

	db, err := sql.Open("postgres", "postgres://danielleknudson:@localhost/gobase?sslmode=disable")

	if err != nil {
		// TODO: Proper error handling
		log.Fatal(err)
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	table := dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	table.ColMap("FirstName").Rename("first_name")
	table.ColMap("LastName").Rename("last_name")
	table.ColMap("CreatedAt").Rename("created_at")
	table.ColMap("ModifiedAt").Rename("modified_at")

	err = dbmap.CreateTablesIfNotExists()

	if err != nil {
		// TODO: Proper error handling
		log.Fatal(err)
	}
}
