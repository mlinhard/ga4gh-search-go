package schema

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"testing"
)

const DEBUG_FILES = false

type MyStruct1 struct {
	Name       string   `json:"name"`
	Age        int      `json:"age"`
	Address    *Address `json:"address,omitempty"`
	privateAge int
}

type Address struct {
	Street string `json:"street"`
	Number int    `json:"number"`
}

type Person interface {
	GetName() string
}

func (s *MyStruct1) GetName() string {
	return s.Name
}

func Test_Generate_Primitive_Int(t *testing.T) {
	var v int = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int")
}

func Test_Generate_Primitive_Bool(t *testing.T) {
	var v bool = true
	assertSchemaGenerationAndValidation(t, v, "primitive_bool")
}

func Test_Generate_Primitive_Float32(t *testing.T) {
	var v float32 = 0.09
	assertSchemaGenerationAndValidation(t, v, "primitive_float32")
}

func Test_Generate_Primitive_Float64(t *testing.T) {
	var v float64 = 0.00009
	assertSchemaGenerationAndValidation(t, v, "primitive_float64")
}

func Test_Generate_Primitive_Int8(t *testing.T) {
	var v int8 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int8")
}

func Test_Generate_Primitive_Int16(t *testing.T) {
	var v int16 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int16")
}

func Test_Generate_Primitive_Int32(t *testing.T) {
	var v int32 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int32")
}

func Test_Generate_Primitive_Int64(t *testing.T) {
	var v int64 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_int64")
}

func Test_Generate_Primitive_Uint(t *testing.T) {
	var v uint = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint")
}

func Test_Generate_Primitive_Uint8(t *testing.T) {
	var v uint8 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint8")
}

func Test_Generate_Primitive_Uint16(t *testing.T) {
	var v uint16 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint16")
}

func Test_Generate_Primitive_Uint32(t *testing.T) {
	var v uint32 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint32")
}

func Test_Generate_Primitive_Uint64(t *testing.T) {
	var v uint64 = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uint64")
}

func Test_Generate_Primitive_UintPtr(t *testing.T) {
	var v uintptr = 1
	assertSchemaGenerationAndValidation(t, v, "primitive_uintptr")
}

func Test_Generate_Primitive_String(t *testing.T) {
	var v string = "Value"
	assertSchemaGenerationAndValidation(t, v, "primitive_string")
}

func Test_Generate_Struct(t *testing.T) {
	var v MyStruct1 = MyStruct1{"John", 12, &Address{"Longstreet", 100}, 13}
	assertSchemaGenerationAndValidation(t, v, "struct1")
}

func Test_Generate_Struct2(t *testing.T) {
	var v MyStruct1 = MyStruct1{"John", 12, nil, 13}
	assertSchemaGenerationAndValidation(t, v, "struct2")
}

func Test_Generate_Slice(t *testing.T) {
	var v []string = []string{"Banana", "Apple", "Cherry"}
	assertSchemaGenerationAndValidation(t, v, "slice")
}

func Test_Generate_Array(t *testing.T) {
	var v [4]string
	v[0] = "Apple"
	v[1] = "Banana"
	v[2] = "Cherry"
	v[3] = "Dumbfruit"
	assertSchemaGenerationAndValidation(t, v, "array")
}

func Test_Generate_Interface(t *testing.T) {
	var v *Person
	assertSchemaGenerationError(t, v, "interfaces not supported: schema.Person")
}

func Test_Generate_Map(t *testing.T) {
	var v map[string]string = map[string]string{"1": "Banana", "2": "Apple", "3": "Cherry"}
	assertSchemaGenerationAndValidation(t, v, "map")
}

// helper functions

func assertSchemaGenerationError(t *testing.T, object interface{}, expectedErrorMessage string) {
	_, err := GenerateSchema(object)
	if err == nil {
		t.Errorf("Schema generation for %v is supposed to return error", object)
	} else {
		assert.Equal(t, expectedErrorMessage, err.Error())
	}
}

func assertSchemaGenerationAndValidation(t *testing.T, object interface{}, debugFilePrefix string) {
	bytes := marshall(t, object)
	if DEBUG_FILES {
		createTestDataDirIfNeeded(t, "testdata")
		writeTestFile(t, bytes, "testdata/"+debugFilePrefix+"_object.json")
	}
	schema, err := GenerateSchema(object)
	if err != nil {
		t.Errorf("Error generating schema for %v", debugFilePrefix)
	}
	schemaBytes := marshall(t, schema)
	if DEBUG_FILES {
		writeTestFile(t, schemaBytes, "testdata/"+debugFilePrefix+"_schema.json")
	}
	validate(t, bytes, schemaBytes)
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
		t.Errorf("Validation error: %v", result)
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
