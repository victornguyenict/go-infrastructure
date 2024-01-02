package jsonutil

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Marshal takes an input Go data structure and returns its JSON encoding
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// PrettyPrint takes a JSON byte slice and prints it in a human-readable format
func PrettyPrint(data []byte) (string, error) {
	var obj map[string]interface{}
	if err := Unmarshal(data, &obj); err != nil {
		return "", err
	}

	prettyJSON, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return "", err
	}

	return string(prettyJSON), nil
}

// ReadJSONFromFile reads a file and unmarshals its content into a Go data structure
func ReadJSONFromFile(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return Unmarshal(data, v)
}

// WriteJSONToFile takes a data structure and writes it as JSON to a file
func WriteJSONToFile(filename string, v interface{}) error {
	data, err := Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, os.ModePerm)
}

// GetJSONFromURL fetches a JSON document from the specified URL and decodes it into the provided variable.
func GetJSONFromURL(url string, v interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(v)
}

// JSONExists checks if a key exists in a given JSON object (only for top-level keys).
func JSONExists(jsonData []byte, key string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonData, &obj); err != nil {
		return false
	}
	_, exists := obj[key]
	return exists
}

// PrettyPrintToWriter takes a JSON byte slice and writes the pretty-printed version to the provided writer.
func PrettyPrintToWriter(data []byte, writer io.Writer) error {
	prettyData, err := PrettyPrint(data)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(prettyData))
	return err
}

// FlattenJSON converts nested JSON object into a flat map with compound keys.
func FlattenJSON(data []byte) (map[string]interface{}, error) {
	var original map[string]interface{}
	err := json.Unmarshal(data, &original)
	if err != nil {
		return nil, err
	}

	flatMap := make(map[string]interface{})
	flatten("", original, flatMap)
	return flatMap, nil
}

// Helper function to recursively flatten the JSON.
func flatten(currentPath string, value interface{}, flatMap map[string]interface{}) {
	if castedMap, ok := value.(map[string]interface{}); ok {
		for k, v := range castedMap {
			newPath := fmt.Sprintf("%s.%s", currentPath, k)
			flatten(newPath, v, flatMap)
		}
	} else {
		trimmedPath := strings.TrimPrefix(currentPath, ".")
		flatMap[trimmedPath] = value
	}
}

// MergeJSON merges two JSON objects, with values from the second object overriding those in the first if they conflict.
func MergeJSON(jsonData1, jsonData2 []byte) ([]byte, error) {
	var obj1, obj2 map[string]interface{}

	err := json.Unmarshal(jsonData1, &obj1)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData2, &obj2)
	if err != nil {
		return nil, err
	}

	for k, v := range obj2 {
		obj1[k] = v
	}

	return json.Marshal(obj1)
}

// ExtractValue gets a value specified by a dot-separated path in a nested JSON structure.
func ExtractValue(jsonData []byte, path string) (interface{}, error) {
	var obj interface{}

	err := json.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}

	pathParts := strings.Split(path, ".")

	for _, part := range pathParts {
		if objMap, ok := obj.(map[string]interface{}); ok {
			obj = objMap[part]
		} else {
			return nil, fmt.Errorf("path not found: %s", path)
		}
	}

	return obj, nil
}

// TransformJSON applies a transformation function to each value in the JSON object.
func TransformJSON(jsonData []byte, transformFunc func(interface{}) interface{}) ([]byte, error) {
	var obj interface{}
	err := json.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}
	transform(obj, transformFunc)
	return json.Marshal(obj)
}

// Helper function to recursively apply the transformation.
func transform(value interface{}, transformFunc func(interface{}) interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		for k, val := range v {
			v[k] = transformFunc(val)
			transform(val, transformFunc)
		}
	case []interface{}:
		for i, val := range v {
			v[i] = transformFunc(val)
			transform(val, transformFunc)
		}
	default:
		value = transformFunc(value)
	}
}

// FilterJSON filters the JSON object based on the provided filter function.
func FilterJSON(jsonData []byte, filterFunc func(string, interface{}) bool) ([]byte, error) {
	var obj map[string]interface{}
	err := json.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}
	for k, v := range obj {
		if !filterFunc(k, v) {
			delete(obj, k)
		}
	}
	return json.Marshal(obj)
}

// DefaultValue sets a default value if the specified key is missing or null in the JSON object.
func DefaultValue(jsonData []byte, key string, defaultValue interface{}) ([]byte, error) {
	var obj map[string]interface{}

	err := json.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}

	if _, ok := obj[key]; !ok || obj[key] == nil {
		obj[key] = defaultValue
	}

	return json.Marshal(obj)
}

// ConvertTypes converts all instances of a type to another within the JSON object.
func ConvertTypes(jsonData []byte, fromType, toType reflect.Kind) ([]byte, error) {
	var obj interface{}
	err := json.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}
	convertTypes(obj, fromType, toType)
	return json.Marshal(obj)
}

// Helper function to recursively convert types.
func convertTypes(value interface{}, fromType, toType reflect.Kind) {
	switch v := value.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if reflect.ValueOf(val).Kind() == fromType {
				v[k] = convertValue(val, toType)
			}
			convertTypes(val, fromType, toType)
		}
	case []interface{}:
		for i, val := range v {
			if reflect.ValueOf(val).Kind() == fromType {
				v[i] = convertValue(val, toType)
			}
			convertTypes(val, fromType, toType)
		}
	}
}

// convertValue converts a single value from one type to another.
func convertValue(value interface{}, toType reflect.Kind) interface{} {
	switch toType {
	case reflect.String:
		return fmt.Sprint(value)
	case reflect.Int, reflect.Int64:
		if val, err := strconv.ParseInt(fmt.Sprint(value), 10, 64); err == nil {
			return val
		}
		// Add more type conversion cases as needed.
	}
	return value
}

// SliceJSONObjects divides a JSON array into smaller slices of a specified size.
func SliceJSONObjects(jsonData []byte, size int) ([][]byte, error) {
	var objs []interface{}
	err := json.Unmarshal(jsonData, &objs)
	if err != nil {
		return nil, err
	}

	var chunks [][]byte

	for i := 0; i < len(objs); i += size {
		end := i + size
		if end > len(objs) {
			end = len(objs)
		}

		chunk, err := json.Marshal(objs[i:end])
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// CompactJSON removes all white spaces and line breaks from a JSON string.
func CompactJSON(jsonData []byte) ([]byte, error) {
	var obj interface{}
	err := json.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}
	return json.Marshal(obj) // json.Marshal removes white spaces and line breaks.
}
