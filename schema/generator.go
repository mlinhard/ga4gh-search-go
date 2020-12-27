package schema

import (
	"fmt"
	"github.com/sourcegraph/go-jsonschema/jsonschema"
	"reflect"
	"strings"
)

func GenerateSchema(v interface{}) (*jsonschema.Schema, error) {
	if v == nil {
		return EmptySchema(), nil
	}
	return generateSchema(reflect.ValueOf(v).Type())
}

func EmptySchema() *jsonschema.Schema {
	return schemaBase(nil, nil, nil)
}

func generateSchema(v reflect.Type) (*jsonschema.Schema, error) {
	switch v.Kind() {
	case reflect.Bool:
		return schemaPrim(jsonschema.BooleanType), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return schemaPrim(jsonschema.IntegerType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return schemaPrim(jsonschema.IntegerType), nil
	case reflect.Float32, reflect.Float64:
		return schemaPrim(jsonschema.NumberType), nil
	case reflect.String:
		return schemaPrim(jsonschema.StringType), nil
	case reflect.Interface:
		return nil, fmt.Errorf("interfaces not supported: %v", v)
	case reflect.Struct:
		return generateStructSchema(v)
	case reflect.Map:
		return schemaBase(jsonschema.PrimitiveTypeList{jsonschema.ObjectType}, nil, nil), nil
	case reflect.Slice:
		return generateArraySchema(v)
	case reflect.Array:
		return generateArraySchema(v)
	case reflect.Ptr:
		return generateSchema(v.Elem())
	default:
		return nil, fmt.Errorf("unsupported type: %v", v)
	}
}

func generateStructSchema(v reflect.Type) (*jsonschema.Schema, error) {
	properties := map[string]*jsonschema.Schema{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		isUnexported := field.PkgPath != ""
		if !isUnexported {
			fieldName := fieldJsonName(field)
			fieldSchema, fieldErr := generateSchema(field.Type)
			if fieldErr != nil {
				return nil, fmt.Errorf("processing field %v of type %v", fieldName, v.Name())
			}
			properties[fieldName] = fieldSchema
		}
	}
	return schemaBase(jsonschema.PrimitiveTypeList{jsonschema.ObjectType}, &properties, nil), nil
}

func generateArraySchema(v reflect.Type) (*jsonschema.Schema, error) {
	elemSchema, err := generateSchema(v.Elem())
	if err != nil {
		return nil, err
	}
	return schemaBase(jsonschema.PrimitiveTypeList{jsonschema.ArrayType}, nil, &jsonschema.SchemaOrSchemaList{
		Schema:  elemSchema,
		Schemas: nil,
	}), nil
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

func schemaPrim(primitive jsonschema.PrimitiveType) *jsonschema.Schema {
	return schemaBase(jsonschema.PrimitiveTypeList{primitive}, nil, nil)
}

// base schema constructor
func schemaBase(typez jsonschema.PrimitiveTypeList, properties *map[string]*jsonschema.Schema, items *jsonschema.SchemaOrSchemaList) *jsonschema.Schema {
	return &jsonschema.Schema{
		Comment:              nil,
		ID:                   nil,
		Reference:            nil,
		SchemaRef:            nil,
		AdditionalItems:      nil,
		AdditionalProperties: nil,
		AllOf:                nil,
		AnyOf:                nil,
		Const:                nil,
		Contains:             nil,
		Default:              nil,
		Definitions:          nil,
		Dependencies:         nil,
		Description:          nil,
		Else:                 nil,
		Enum:                 nil,
		Examples:             nil,
		ExclusiveMaximum:     nil,
		ExclusiveMinimum:     nil,
		Format:               nil,
		If:                   nil,
		Items:                items,
		MaxItems:             nil,
		MaxLength:            nil,
		MaxProperties:        nil,
		Maximum:              nil,
		MinItems:             nil,
		MinLength:            nil,
		MinProperties:        nil,
		Minimum:              nil,
		MultipleOf:           nil,
		Not:                  nil,
		OneOf:                nil,
		Pattern:              nil,
		PatternProperties:    nil,
		Properties:           properties,
		PropertyNames:        nil,
		Required:             nil,
		Then:                 nil,
		Title:                nil,
		Type:                 typez,
		UniqueItems:          nil,
		Raw:                  nil,
		IsEmpty:              false,
		IsNegated:            false,
		Go:                   nil,
	}
}
