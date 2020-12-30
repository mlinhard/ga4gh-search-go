# GA4GH Search Backend

This is an experimental implementation of [GA4GH-Search](https://github.com/ga4gh-discovery/ga4gh-search)
[API](https://github.com/ga4gh-discovery/ga4gh-search/blob/develop/openapi/openapi.yaml) in Go

## Dependencies

So far the code depends on these third-party non-core libraries

- [sourcegraph/go-jsonschema](https://github.com/sourcegraph/go-jsonschema) - structs
  for [JSON Schema](https://json-schema.org/) model serialization
- [gorilla/mux](https://github.com/gorilla/mux) - http server router
- [xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema) - for [JSON Schema](https://json-schema.org/)
  validation in tests
- [stretchr/testify](https://github.com/stretchr/testify) - for assertions in tests

