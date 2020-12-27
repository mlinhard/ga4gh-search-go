package api

import (
	"github.com/sourcegraph/go-jsonschema/jsonschema"
)

// Go struct encoding of GA4GH Search API
// OpenAPI v3 Spec: https://github.com/ga4gh-discovery/ga4gh-search/blob/develop/openapi/openapi.yaml

type Search interface {
	// GET /tables
	Tables() (*ListTablesResponse, error)

	// GET /table/{table_name}/info
	TableInfo(name string) (*Table, error)

	// GET /table/{table_name}/data
	TableData(name string) (*TableData, error)

	// POST /search
	Search(request *SearchRequest) (*TableData, error)

	// GET /service-info
	ServiceInfo() (*Service, error)
}

type ListTablesResponse struct {
	Tables     []*Table    `json:"tables"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Errors     ErrorList   `json:"errors,omitempty"`
}

type Table struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	DataModel   *jsonschema.Schema `json:"data_model"`
	Errors      ErrorList          `json:"errors,omitempty"`
}

type TableData struct {
	DataModel  *jsonschema.Schema `json:"data_model"`
	Data       []interface{}      `json:"data"`
	Pagination *Pagination        `json:"pagination,omitempty"`
	Errors     ErrorList          `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Errors ErrorList `json:"errors"`
}

type ErrorList []*Error

type Error struct {
	Source string `json:"source"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type Pagination struct {
	NextPageUrl string `json:"next_page_url"`
}
type SearchRequest struct {
	Query      string        `json:"query"`
	Parameters []interface{} `json:"parameters,omitempty"`
}

// Go struct encoding of GA4GH service-info API
// OpenAPI v3 Spec: https://github.com/ga4gh-discovery/ga4gh-service-info/blob/develop/service-info.yaml

type Service struct {
	Id               string       `json:"id"`
	Name             string       `json:"name"`
	Type             *ServiceType `json:"type"`
	Description      string       `json:"description"`
	Organization     string       `json:"organization"`
	ContactUrl       string       `json:"contactUrl"`
	DocumentationUrl string       `json:"documentationUrl"`
	CreatedAt        string       `json:"createdAt"`
	UpdatedAt        string       `json:"updatedAt"`
	Environment      string       `json:"environment"`
	Version          string       `json:"version"`
}

type ServiceType struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
}
