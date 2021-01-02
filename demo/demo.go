package main

import (
	"fmt"
	"github.com/mlinhard/ga4gh-search-go/schema"
	"github.com/mlinhard/ga4gh-search-go/server"
	"github.com/mlinhard/ga4gh-search-go/sql"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Sex string
type JSONDate time.Time

const (
	M Sex = "M"
	F     = "F"
	X     = "X"
)

type Person struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Born    JSONDate `json:"born"`
	Sex     Sex      `json:"sex"`
}

func main() {
	sql.Main()
}

func main2() {
	// create HTTP server
	server, err := server.NewHttpSearchServer("localhost:8080")
	if err != nil {
		fmt.Printf("Error starting server: %v", err)
		return
	}

	// create data
	persons := []interface{}{
		Person{
			Name:    "Alfred",
			Surname: "Adams",
			Sex:     "M",
			Born:    date("2020-01-01"),
		}, Person{
			Name:    "Betty",
			Surname: "Bumst",
			Sex:     "F",
			Born:    date("2020-02-02"),
		}, Person{
			Name:    "Cecille",
			Surname: "Chang",
			Sex:     "M",
			Born:    date("2020-03-03"),
		}}

	// load data schema
	personSchema, err := schema.LoadSchema("person_schema.json")
	if err != nil {
		fmt.Printf("ERROR loading schema: %v", err)
	}

	search := server.Service()

	// create table
	search.AddTable("persons", "Table of persons", persons, personSchema)

	// add enumerations so that proper schema is generated for enums
	search.SchemaHintEnum(M, F, X)
	search.SchemaHintFormat(JSONDate(time.Now()), "date")

	// create table with auto-generated schema
	err = search.AddTableAutoSchema("persons_auto", "Table of persons with auto-generated schema", persons)
	if err != nil {
		fmt.Printf("ERROR generating schema: %v", err)
	}

	// start the server and wait for termination signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	server.Start()
	<-done
	server.Stop()
	fmt.Printf("Server terminated.\n")
}

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return ([]byte)("\"" + time.Time(d).Format("2006-01-02") + "\""), nil
}

func date(date string) JSONDate {
	v, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return JSONDate(v)
}
