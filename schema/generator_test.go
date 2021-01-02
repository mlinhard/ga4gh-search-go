package schema

import (
	"encoding/json"
	"fmt"
	. "github.com/mlinhard/ga4gh-search-go/tests"
	"github.com/sourcegraph/go-jsonschema/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

const DEBUG_FILES = false

func Test_Generate_Primitive_Int(t *testing.T) {
	var v int = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Bool(t *testing.T) {
	var v bool = true
	assertSchemaGenerationAndValidation(t, v, "primitive_bool").
		AssertType(jsonschema.BooleanType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Float32(t *testing.T) {
	var v float32 = 0.09
	assertSchemaGenerationAndValidation(t, v, "primitive_float32").
		AssertType(jsonschema.NumberType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Float64(t *testing.T) {
	var v float64 = 0.00009
	assertSchemaGenerationAndValidation(t, v, "primitive_float64").
		AssertType(jsonschema.NumberType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Int8(t *testing.T) {
	var v int8 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int8").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Int16(t *testing.T) {
	var v int16 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int16").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Int32(t *testing.T) {
	var v int32 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int32").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Int64(t *testing.T) {
	var v int64 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int64").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Uint(t *testing.T) {
	var v uint = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Uint8(t *testing.T) {
	var v uint8 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint8").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Uint16(t *testing.T) {
	var v uint16 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint16").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Uint32(t *testing.T) {
	var v uint32 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint32").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_Uint64(t *testing.T) {
	var v uint64 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint64").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_UintPtr(t *testing.T) {
	var v uintptr = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uintptr").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Primitive_String(t *testing.T) {
	var v string = "Value"
	assertSchemaGenerationAndValidation(t, v, "primitive_string").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
}

func Test_Generate_Struct(t *testing.T) {
	var v Person = Person{"John", 12, M, &Address{"Longstreet", 100}, time.Now(), date("2020-01-01"), 1}

	schema := assertSchemaGenerationAndValidation(t, v, "struct1").
		AssertType(jsonschema.ObjectType).
		AssertNumProperties(6).
		AssertOnlyDef("Type", "Properties")

	schema.Property("name").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
	schema.Property("age").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")

	addrSchema := schema.Property("address").
		AssertType(jsonschema.ObjectType).
		AssertOnlyDef("Type", "Properties")
	addrSchema.Property("street").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
	addrSchema.Property("number").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")

	schema.Property("updated_at").
		AssertType(jsonschema.StringType).
		AssertFormat("date-time").
		AssertOnlyDef("Type", "Format")
	schema.Property("born").
		AssertType(jsonschema.StringType).
		AssertFormat("date").
		AssertOnlyDef("Type", "Format")
	schema.Property("sex").
		AssertType(jsonschema.StringType).
		AssertEnums(M, F, X).
		AssertOnlyDef("Type", "Enum")
}

func Test_Generate_Struct2(t *testing.T) {
	var v Person = Person{"John", 12, M, nil, time.Now(), date("2020-01-01"), 13}
	// schema will be generated for address sub-schema even when the value is nil
	schema := assertSchemaGenerationAndValidation(t, v, "struct2")
	addrSchema := schema.Property("address").
		AssertType(jsonschema.ObjectType).
		AssertOnlyDef("Type", "Properties")
	addrSchema.Property("street").
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
	addrSchema.Property("number").
		AssertType(jsonschema.IntegerType).
		AssertOnlyDef("Type")
}

func Test_Generate_Slice(t *testing.T) {
	var v []string = []string{"Banana", "Apple", "Cherry"}
	assertSchemaGenerationAndValidation(t, v, "slice").
		AssertType(jsonschema.ArrayType).
		AssertOnlyDef("Type", "Items").
		Items().
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
}

func Test_Generate_Array(t *testing.T) {
	var v [4]string
	v[0] = "Apple"
	v[1] = "Banana"
	v[2] = "Cherry"
	v[3] = "Dumbfruit"
	assertSchemaGenerationAndValidation(t, v, "array").
		AssertType(jsonschema.ArrayType).
		AssertOnlyDef("Type", "Items").
		Items().
		AssertType(jsonschema.StringType).
		AssertOnlyDef("Type")
}

func Test_Generate_Interface(t *testing.T) {
	var v *IPerson
	assertSchemaGenerationError(t, v, "interfaces not supported: schema.IPerson")
}

func Test_Generate_Map(t *testing.T) {
	var v map[string]string = map[string]string{"1": "Banana", "2": "Apple", "3": "Cherry"}
	assertSchemaGenerationAndValidation(t, v, "map").
		AssertType(jsonschema.ObjectType).
		AssertOnlyDef("Type")
}

// helper functions

func assertSchemaGenerationError(t *testing.T, object interface{}, expectedErrorMessage string) {
	_, err := NewGenerator().GenerateSchema(object)
	if err == nil {
		t.Errorf("Schema generation for %v is supposed to return error", object)
	} else {
		assert.Equal(t, expectedErrorMessage, err.Error())
	}
}

func assertSchemaGenerationAndValidation(t *testing.T, object interface{}, debugFilePrefix string) *SchemaTester {
	bytes := marshall(t, object)
	if DEBUG_FILES {
		createTestDataDirIfNeeded(t, "debug_schemas")
		writeTestFile(t, bytes, "debug_schemas/"+debugFilePrefix+"_object.json")
	}
	gen := NewGenerator()
	gen.SchemaHintEnum(reflect.TypeOf(M), jsonschema.EnumList{M, F, X})
	fmt := jsonschema.Format("date")
	gen.SchemaHintFormat(reflect.TypeOf(JSONDate(time.Time{})), &fmt)
	schema, err := gen.GenerateSchema(object)
	if err != nil {
		t.Errorf("Error generating schema for %v", debugFilePrefix)
	}
	schemaBytes := marshall(t, schema)
	if DEBUG_FILES {
		writeTestFile(t, schemaBytes, "debug_schemas/"+debugFilePrefix+"_schema.json")
	}
	validate(t, bytes, schemaBytes)
	return NewSchemaTester(t, schema)
}

func validate(t *testing.T, jsonBytes []byte, schemaBytes []byte) {
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Errorf("Error validating schema: %v", err)
		return
	}
	if !result.Valid() {
		msg := fmt.Sprintf("DOCUMENT:\n%v\nSCHEMA:\n%v\nVALIDATION ERRORS:\n", string(jsonBytes), string(schemaBytes))
		for _, error := range result.Errors() {
			msg += error.String() + "\n"
		}
		t.Errorf(msg)
	}
}

func marshall(t *testing.T, object interface{}) []byte {
	b, err := json.Marshal(object)
	if err != nil {
		t.Errorf("Couldn't marshall")
	}
	return b
}

func writeTestFile(t *testing.T, bytes []byte, filename string) {
	err := ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		t.Errorf("Couldn't write")
	}
}

func createTestDataDirIfNeeded(t *testing.T, dir string) {
	_, errStat := os.Stat(dir)
	if os.IsNotExist(errStat) {
		errMkdir := os.Mkdir(dir, os.ModePerm)
		if errMkdir != nil {
			t.Errorf("Couldn't create directory %v", dir)
		}
	}
}

type Sex string

const (
	M Sex = "M"
	F     = "F"
	X     = "X"
)

type JSONDate time.Time

type Person struct {
	Name       string    `json:"name"`
	Age        int       `json:"age"`
	Sex        Sex       `json:"sex"`
	Address    *Address  `json:"address,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
	Born       JSONDate  `json:"born"`
	privateAge int
}

type Address struct {
	Street string `json:"street"`
	Number int    `json:"number"`
}

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return ([]byte)("\"" + time.Time(d).Format("2006-01-02") + "\""), nil
}

func date(date string) JSONDate {
	v, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return JSONDate(v)
}

type IPerson interface {
	GetName() string
}

func (s *Person) GetName() string {
	return s.Name
}
