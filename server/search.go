package server

import (
	"github.com/mlinhard/ga4gh-search-go/api"
	"github.com/mlinhard/ga4gh-search-go/schema"
	"github.com/sourcegraph/go-jsonschema/jsonschema"
	"reflect"
)

type table struct {
	name        string
	description string
	schema      *jsonschema.Schema
	data        []interface{}
}

// this is an internal server api and will be subject to heavy changes in initial stages of development
type SearchService struct {
	tables          map[string]*table
	schemaGenerator *schema.Generator
}

func NewSearchService() (*SearchService, error) {
	server := new(SearchService)
	server.tables = make(map[string]*table)
	server.schemaGenerator = schema.NewGenerator()
	return server, nil
}

func (s *SearchService) SchemaHintEnum(values ...interface{}) {
	if len(values) == 0 {
		panic("empty enumeration")
	}
	s.schemaGenerator.SchemaHintEnum(reflect.TypeOf(values[0]), jsonschema.EnumList{values})
}

func (s *SearchService) SchemaHintFormat(v interface{}, format string) {
	var jsonSchemaFormat = jsonschema.Format(format)
	s.schemaGenerator.SchemaHintFormat(reflect.TypeOf(v), &jsonSchemaFormat)
}

func (s *SearchService) SchemaHintPattern(v reflect.Type, pattern *string) {
	s.schemaGenerator.SchemaHintPattern(v, pattern)
}

func (s *SearchService) AddTableAutoSchema(name string, description string, data []interface{}) error {
	schema, err := s.generateDataModel(data)
	if err != nil {
		return err
	}
	s.tables[name] = &table{name: name, description: description, data: data, schema: schema}
	return nil
}

func (s *SearchService) AddTable(name string, description string, data []interface{}, schema *jsonschema.Schema) {
	s.tables[name] = &table{name: name, description: description, data: data, schema: schema}
}

func (s *SearchService) Tables() (*api.ListTablesResponse, error) {
	var result = make([]*api.Table, 0, len(s.tables))
	for _, table := range s.tables {
		apiTable, err := s.toApi(table)
		if err != nil {
			return nil, err
		}
		result = append(result, apiTable)
	}
	return &api.ListTablesResponse{
		Tables:     result,
		Pagination: nil,
		Errors:     nil,
	}, nil
}

func (s *SearchService) toApi(table *table) (*api.Table, error) {
	return &api.Table{
		Name:        table.name,
		Description: table.description,
		DataModel:   table.schema,
		Errors:      nil,
	}, nil
}

func (s *SearchService) generateDataModel(data []interface{}) (*jsonschema.Schema, error) {
	var v interface{}
	if data == nil || len(data) == 0 {
		v = nil
	} else {
		v = data[0]
	}
	return s.schemaGenerator.GenerateSchema(v)
}

// GET /table/{table_name}/info
func (s *SearchService) TableInfo(name string) (*api.Table, error) {
	table := s.tables[name]
	if table == nil {
		return nil, nil
	}
	return s.toApi(table)
}

// GET /table/{table_name}/data
func (s *SearchService) TableData(name string) (*api.TableData, error) {
	table := s.tables[name]
	if table == nil {
		return nil, nil
	}
	apiTable, err := s.toApi(table)
	if err != nil {
		return nil, err
	}
	return &api.TableData{
		DataModel:  apiTable.DataModel,
		Data:       table.data,
		Pagination: nil,
		Errors:     nil,
	}, nil
}

// POST /search
func (s *SearchService) Search(request *api.SearchRequest) (*api.TableData, error) {
	return &api.TableData{
		DataModel:  schema.EmptySchema(),
		Data:       []interface{}{},
		Pagination: nil,
		Errors:     nil,
	}, nil
}

// GET /service-info
func (s *SearchService) ServiceInfo() (*api.Service, error) {
	return &api.Service{
		Id:   "search-service",
		Name: "GA4GH Search Service",
		Type: &api.ServiceType{
			Group:    "sk.linhard.search",
			Artifact: "ga4gh-search",
			Version:  "0.1.0",
		},
		Description:      "GA4GH Search Service",
		Organization:     "Linhard, s.r.o.",
		ContactUrl:       "https://sro.linhard.sk",
		DocumentationUrl: "https://sro.linhard.sk",
		CreatedAt:        "2020-12-29",
		UpdatedAt:        "2020-12-29",
		Environment:      "test",
		Version:          "0.1.0",
	}, nil
}
