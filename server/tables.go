package server

import (
	"encoding/json"
	. "github.com/mlinhard/ga4gh-search-go/api"
	"github.com/sourcegraph/go-jsonschema/jsonschema"
	"io"
	"os"
)

type TableRepo struct {
	tables []Table
}

func NewTableRepo() (*TableRepo, error) {
	tr := new(TableRepo)
	tr.addDemoTable()
	return tr, nil
}

func (t *TableRepo) add(table Table) {
	t.tables = append(t.tables, table)
}

func (t *TableRepo) getTables() []Table {
	return t.tables
}
func (t *TableRepo) getTable(name string) *Table {
	for _, table := range t.tables {
		if table.Name == name {
			return &table
		}
	}
	return nil
}

func (t *TableRepo) getTableData(name string) *TableData {
	table := t.getTable(name)
	return &TableData{
		DataModel: table.DataModel,
		Data:      []interface{}{},
	}
}

func (t *TableRepo) addDemoTable() {
	personSchema, err := readSchema("main/testdata/person_schema.json")
	if err != nil {
		panic(err)
	}
	t.add(Table{
		DataModel:   personSchema,
		Name:        "persons",
		Description: "Table of persons",
	})
}

func readSchema(filename string) (*jsonschema.Schema, error) {
	var f io.ReadCloser
	if filename == "-" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()

	var schema *jsonschema.Schema
	if err := json.NewDecoder(f).Decode(&schema); err != nil {
		return nil, err
	}
	return schema, nil
}
