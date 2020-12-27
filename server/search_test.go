package server

import (
	"github.com/mlinhard/ga4gh-search-go/api"
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
	server.AddTable("persons", "Table of persons", getDemoPersons())
	server.AddTable("pets", "Table of pets", getDemoPets())

	response, err := server.Tables()
	assert.Equal(t, 2, len(response.Tables))

	tables := toTableMap(response.Tables)

	personsTable := tables["persons"]
	assert.Equal(t, "persons", personsTable.Name)
	assert.Equal(t, "Table of persons", personsTable.Description)

	petsTable := tables["pets"]
	assert.Equal(t, "pets", petsTable.Name)
	assert.Equal(t, "Table of pets", petsTable.Description)

	personSchema := personsTable.DataModel
	assertType(t, jsonschema.ObjectType, personSchema)

	assert.Equal(t, 4, len(*personSchema.Properties))

	assertType(t, jsonschema.StringType, (*personSchema.Properties)["name"])
	assertType(t, jsonschema.StringType, (*personSchema.Properties)["surname"])
	assertType(t, jsonschema.StringType, (*personSchema.Properties)["sex"])
	assertType(t, jsonschema.StringType, (*personSchema.Properties)["born"])
}

func assertType(t *testing.T, expectedType jsonschema.PrimitiveType, schema *jsonschema.Schema) {
	assert.Equal(t, 1, schema.Type.Len())
	assert.Equal(t, expectedType, schema.Type[0])
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
