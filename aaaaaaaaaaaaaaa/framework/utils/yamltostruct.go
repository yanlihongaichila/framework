package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

// Helper function to convert map[interface{}]interface{} to map[string]interface{}
func convertMap(input map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range input {
		strKey, ok := key.(string)
		if !ok {
			continue // or log an error, as you prefer
		}

		// Convert nested maps
		if nestedMap, ok := value.(map[interface{}]interface{}); ok {
			value = convertMap(nestedMap)
		}
		result[strKey] = value
	}
	return result
}

// Updated fillStruct function
func FillStruct(v reflect.Value, data interface{}) error {
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("expected a non-nil pointer to a struct")
	}

	v = v.Elem() // Dereference the pointer to get the struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct")
	}

	// Convert data to the correct type
	var dataMap map[string]interface{}
	switch dataTyped := data.(type) {
	case map[string]interface{}:
		dataMap = dataTyped
	case map[interface{}]interface{}:
		dataMap = convertMap(dataTyped)
	default:
		return fmt.Errorf("unexpected data type")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		// Get the YAML tag for the field, fallback to field name if no tag
		yamlKey := fieldType.Tag.Get("yaml")
		if yamlKey == "" {
			yamlKey = strings.ToLower(fieldType.Name[:1]) + fieldType.Name[1:]
		}

		if fieldValue, ok := dataMap[yamlKey]; ok {
			switch field.Kind() {
			case reflect.String:
				if strVal, ok := fieldValue.(string); ok {
					field.SetString(strVal)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// Convert numeric values as necessary
				field.SetInt(convertToInt64(fieldValue))
			case reflect.Struct:
				if nestedMap, ok := fieldValue.(map[string]interface{}); ok {
					if err := FillStruct(field.Addr(), nestedMap); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// Utility function to safely convert interface types to int64, used for handling YAML numeric values.
func convertToInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v) // YAML unmarshalling may treat numbers as floats
	default:
		return 0 // Could log an error or return an error if needed
	}
}

func UnmarshalYaml(structPtr interface{}, data string) error {
	var yamlData map[interface{}]interface{}
	if err := yaml.Unmarshal([]byte(data), &yamlData); err != nil {
		fmt.Printf("Error parsing YAML data: %s\n", err)
		return err
	}

	if err := FillStruct(reflect.ValueOf(structPtr), yamlData); err != nil {
		fmt.Printf("Error filling struct: %s\n", err)
		return err
	}
	return nil
}
