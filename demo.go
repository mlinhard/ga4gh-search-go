package main

import (
	"fmt"
	"github.com/mlinhard/ga4gh-search-go/server"
	"os"
	"os/signal"
	"syscall"
)

type Person struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Sex     string `json:"sex"`
	Born    string `json:"born"`
}

func main() {
	server, err := server.NewHttpSearchServer("localhost:8080")
	if err != nil {
		fmt.Printf("Error starting server: %v", err)
		return
	}

	server.Service().AddTable("persons", "Table of persons", getDemoPersons())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	server.Start()
	<-done
	server.Stop()
	fmt.Printf("Server terminated.\n")
}

func getDemoPersons() []interface{} {
	return []interface{}{
		Person{
			Name:    "John",
			Surname: "Adams",
			Sex:     "M",
			Born:    "1980-01-01",
		}, Person{
			Name:    "Betty",
			Surname: "Bumst",
			Sex:     "F",
			Born:    "1980-02-02",
		}}
}
