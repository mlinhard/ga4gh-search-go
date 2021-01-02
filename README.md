# GA4GH Search Backend

This is an experimental implementation of [GA4GH-Search](https://github.com/ga4gh-discovery/ga4gh-search)
[API](https://github.com/ga4gh-discovery/ga4gh-search/blob/develop/openapi/openapi.yaml) in Go

## Dependencies

So far the code depends on these third-party non-core libraries

- [SqlBase.g4](sql/SqlBase.g4) (SQL grammar) based on [SqlBase.g4](https://github.com/trinodb/trino/blob/master/presto-parser/src/main/antlr4/io/prestosql/sql/parser/SqlBase.g4)
  in [Trino DB](https://github.com/trinodb/trino) project
- [sourcegraph/go-jsonschema](https://github.com/sourcegraph/go-jsonschema) - structs
  for [JSON Schema](https://json-schema.org/) model serialization
- [gorilla/mux](https://github.com/gorilla/mux) - http server router
- [xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema) - for [JSON Schema](https://json-schema.org/)
  validation in tests
- [stretchr/testify](https://github.com/stretchr/testify) - for assertions in tests

