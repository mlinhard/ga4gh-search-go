package server

import (
	"github.com/mlinhard/ga4gh-search-go/api"
	. "github.com/mlinhard/ga4gh-search-go/tests"
	"github.com/sourcegraph/go-jsonschema/jsonschema"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Sex     string `json:"sex"`
	Born    string `json:"born"`
}

type Pet struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func TestSearch(t *testing.T) {
	server, err := NewSearchService()
	if err != nil {
		t.Errorf("creating server: %v", err)
	}
	server.AddTableAutoSchema("persons", "Table of persons", getDemoPersons())
	server.AddTableAutoSchema("pets", "Table of pets", getDemoPets())

	response, err := server.Tables()
	assert.Equal(t, 2, len(response.Tables))

	tables := toTableMap(response.Tables)

	personsTable := tables["persons"]
	assert.Equal(t, "persons", personsTable.Name)
	assert.Equal(t, "Table of persons", personsTable.Description)

	petsTable := tables["pets"]
	assert.Equal(t, "pets", petsTable.Name)
	assert.Equal(t, "Table of pets", petsTable.Description)

	personSchema := schemaTester(t, personsTable).
		AssertType(jsonschema.ObjectType).
		AssertNumProperties(4).
		AssertOnlyDef("Type", "Properties")

	personSchema.Property("name").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
	personSchema.Property("surname").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
	personSchema.Property("sex").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
	personSchema.Property("born").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
}

func schemaTester(t *testing.T, table *api.Table) *SchemaTester {
	return NewSchemaTester(t, table.DataModel)
}

func toTableMap(tables []*api.Table) map[string]*api.Table {
	tableMap := make(map[string]*api.Table)
	for _, table := range tables {
		tableMap[table.Name] = table
	}
	return tableMap
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

func getDemoPets() []interface{} {
	return []interface{}{
		Pet{
			Name: "Bob",
			Kind: "dog",
		}, Pet{
			Name: "Gloria",
			Kind: "cat",
		}}
}
