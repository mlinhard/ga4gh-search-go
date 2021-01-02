package tests

import (
	"github.com/sourcegraph/go-jsonschema/jsonschema"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type SchemaTester struct {
	Schema *jsonschema.Schema
	t      *testing.T
}

func NewSchemaTester(t *testing.T, schema *jsonschema.Schema) *SchemaTester {
	tester := new(SchemaTester)
	tester.Schema = schema
	tester.t = t
	return tester
}

func (s *SchemaTester) AssertType(expectedType jsonschema.PrimitiveType) *SchemaTester {
	assert.Equal(s.t, 1, s.Schema.Type.Len())
	assert.Equal(s.t, expectedType, s.Schema.Type[0])
	return s
}

func (s *SchemaTester) Property(name string) *SchemaTester {
	properties := *s.Schema.Properties
	assert.NotNil(s.t, properties, "schema doesn't have properties")
	propertySchema := properties[name]
	assert.NotNil(s.t, propertySchema, "property %v not found", name)
	return &SchemaTester{propertySchema, s.t}
}

func (s *SchemaTester) Items() *SchemaTester {
	items := *s.Schema.Items
	assert.NotNil(s.t, items, "schema doesn't have items sub-schema")
	assert.NotNil(s.t, items.Schema, "schema doesn't have items sub-schema")
	return &SchemaTester{items.Schema, s.t}
}

func (s *SchemaTester) ItemsMany(i int) *SchemaTester {
	items := *s.Schema.Items
	assert.NotNil(s.t, items, "schema doesn't have items sub-schema")
	assert.NotNil(s.t, items.Schemas, "schema doesn't have items array sub-schema")
	return &SchemaTester{items.Schemas[i], s.t}
}

func (s *SchemaTester) AssertEnums(values ...interface{}) *SchemaTester {
	assert.NotNil(s.t, s.Schema.Enum, "schema doesn't have enum")
	assert.Equal(s.t, len(values), len(s.Schema.Enum))
	for i, v := range values {
		assert.Equal(s.t, v, s.Schema.Enum[i], "The %v-th value doesn't match", i)
	}
	return s
}

func (s *SchemaTester) AssertOnlyDef(definedProperties ...string) *SchemaTester {
	AssertOnlyDef(s.t, reflect.ValueOf(s.Schema), definedProperties...)
	return s
}

func AssertOnlyDef(t *testing.T, sv reflect.Value, definedProperties ...string) {
	for sv.Type().Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	allowed := make(map[string]bool)
	for _, prop := range definedProperties {
		allowed[prop] = true
	}
	st := sv.Type()
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		fv := sv.Field(i)
		if allowed[ft.Name] {
			assert.False(t, fv.IsZero(), "Field %v should be defined", ft.Name)
		} else {
			assert.True(t, fv.IsZero(), "Field %v should be undefined", ft.Name)
		}
	}
}

func (s *SchemaTester) AssertFormat(expectedFormat string) *SchemaTester {
	assert.NotNil(s.t, *s.Schema.Format, "schema doesn't have format")
	assert.Equal(s.t, expectedFormat, string(*s.Schema.Format))
	return s
}

func (s *SchemaTester) AssertNumProperties(expectedNumProperties int) *SchemaTester {
	properties := *s.Schema.Properties
	assert.NotNil(s.t, properties, "schema doesn't have properties")
	assert.Equal(s.t, expectedNumProperties, len(properties))
	return s
}
