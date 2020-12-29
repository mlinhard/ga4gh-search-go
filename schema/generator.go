package schema

import (
	"encoding/json"
	"fmt"
	"github.com/sourcegraph/go-jsonschema/jsonschema"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

type Generator struct {
	timeType    reflect.Type
	schemaHints map[reflect.Type]*schemaHint
}

type schemaHint struct {
	enum      jsonschema.EnumList
	format    *jsonschema.Format
	pattern   *string
	minLength *int64
	maxLength *int64
}

func LoadSchema(filename string) (*jsonschema.Schema, error) {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(filename); err != nil {
		return nil, err
	}
	schema := new(jsonschema.Schema)
	if err = json.Unmarshal(data, schema); err != nil {
		return nil, err
	}
	return schema, nil
}

func NewGenerator() *Generator {
	g := new(Generator)
	g.timeType = reflect.TypeOf(time.Time{})
	g.schemaHints = make(map[reflect.Type]*schemaHint)
	return g
}

func (g *Generator) SchemaHintEnum(v reflect.Type, values jsonschema.EnumList) {
	g.ensureHint(v).enum = values
}

func (g *Generator) SchemaHintFormat(v reflect.Type, format *jsonschema.Format) {
	g.ensureHint(v).format = format
}

func (g *Generator) SchemaHintPattern(v reflect.Type, pattern *string) {
	g.ensureHint(v).pattern = pattern
}

func (g *Generator) ensureHint(v reflect.Type) *schemaHint {
	hint := g.schemaHints[v]
	if hint == nil {
		hint = new(schemaHint)
		g.schemaHints[v] = hint
	}
	return hint
}

func (g *Generator) GenerateSchema(v interface{}) (*jsonschema.Schema, error) {
	if v == nil {
		return EmptySchema(), nil
	}
	return g.generateSchema(reflect.ValueOf(v).Type())
}

func EmptySchema() *jsonschema.Schema {
	return &jsonschema.Schema{}
}

func (g *Generator) generateSchema(v reflect.Type) (*jsonschema.Schema, error) {
	if v.Kind() == reflect.Ptr {
		return g.generateSchema(v.Elem())
	} else if v.Kind() == reflect.Interface {
		return nil, fmt.Errorf("interfaces not supported: %v", v)
	}
	typeList := g.typeList(v)
	schema := &jsonschema.Schema{Type: typeList}

	// apply specific stuff
	switch typeList[0] {
	case jsonschema.StringType:
		if v.AssignableTo(g.timeType) {
			fmt := jsonschema.Format("date-time")
			schema.Format = &fmt
		}
	case jsonschema.ObjectType:
		props, err := g.generateProperties(v)
		if err != nil {
			return nil, err
		}
		schema.Properties = props
	case jsonschema.ArrayType:
		elemSchema, err := g.generateSchema(v.Elem())
		if err != nil {
			return nil, err
		}
		schema.Items = &jsonschema.SchemaOrSchemaList{
			Schema:  elemSchema,
			Schemas: nil,
		}
	}

	hint := g.schemaHints[v]

	if hint != nil {
		if hint.enum != nil {
			schema.Enum = hint.enum
		}
		if hint.format != nil {
			schema.Format = hint.format
		}
		if hint.pattern != nil {
			schema.Pattern = hint.pattern
		}
		if hint.minLength != nil {
			schema.MinLength = hint.minLength
		}
		if hint.maxLength != nil {
			schema.MaxLength = hint.maxLength
		}
	}

	return schema, nil
}

func (g *Generator) typeList(t reflect.Type) jsonschema.PrimitiveTypeList {
	switch t.Kind() {
	case reflect.Bool:
		return jsonschema.PrimitiveTypeList{jsonschema.BooleanType}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return jsonschema.PrimitiveTypeList{jsonschema.IntegerType}
	case reflect.Float32, reflect.Float64:
		return jsonschema.PrimitiveTypeList{jsonschema.NumberType}
	case reflect.String:
		return jsonschema.PrimitiveTypeList{jsonschema.StringType}
	case reflect.Struct, reflect.Map:
		if t.ConvertibleTo(g.timeType) {
			return jsonschema.PrimitiveTypeList{jsonschema.StringType}
		}
		return jsonschema.PrimitiveTypeList{jsonschema.ObjectType}
	case reflect.Slice, reflect.Array:
		return jsonschema.PrimitiveTypeList{jsonschema.ArrayType}
	default:
		panic("unsupported kind: " + t.Kind().String())
	}
}

func (g *Generator) generateProperties(v reflect.Type) (*map[string]*jsonschema.Schema, error) {
	if v.Kind() == reflect.Map {
		return nil, nil
	}
	properties := map[string]*jsonschema.Schema{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		isUnexported := field.PkgPath != ""
		if !isUnexported {
			fieldName := fieldJsonName(field)
			fieldSchema, fieldErr := g.generateSchema(field.Type)
			if fieldErr != nil {
				return nil, fmt.Errorf("processing field %v of type %v", fieldName, v.Name())
			}
			properties[fieldName] = fieldSchema
		}
	}
	return &properties, nil
}

func fieldJsonName(field reflect.StructField) string {
	if jsonTag, ok := field.Tag.Lookup("json"); ok {
		commaIdx := strings.Index(jsonTag, ",")
		if commaIdx != -1 {
			return jsonTag[:commaIdx]
		}
		return jsonTag
	}
	return field.Name
}
