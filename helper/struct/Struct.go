package helper

import (
	"reflect"
	"strings"
)

type OptionInterface interface {
	GetType() string
	GetValue() interface{}
}

type Prefix struct {
	value string
}

func (p *Prefix) GetType() string {
	return "prefix"
}

func (p *Prefix) GetValue() interface{} {
	return p.value
}

func WithPrefix(prefix string) *Prefix {
	return &Prefix{value: prefix}
}

func MapAsGorm(in interface{}, options ...OptionInterface) map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// @todo Nest down into structs
	if v.Kind() != reflect.Struct {
		return nil
	}

	reflectType := v.Type()

	for i := 0; i < v.NumField(); i++ {

		// Gets us a StructField
		field := reflectType.Field(i)

		// Get the gorm tag
		if tagValue := field.Tag.Get("gorm"); tagValue != "" {

			// Split the Gorm DB tag by any semi-colons
			tagParts := strings.Split(tagValue, ";")

			// Loop over every separated part in the Gorm string, they can be in any order
			for _, tagPart := range tagParts {

				// Now split them by their sub-parts
				innerTagParts := strings.Split(tagPart, ":")

				if len(innerTagParts) < 2 {
					continue
				}

				// We are expected "column:X" whereby we want X
				if innerTagParts[0] == "column" {
					out[innerTagParts[1]] = v.Field(i).Interface()
				}
			}
		}
	}

	// Loop over options, apply the modifications
	for _, option := range options {
		switch option.GetType() {
		case "prefix":
			out = applyPrefix(out, option.GetValue().(string))
		}
	}

	return out
}

func Combine(interfaces ...map[string]interface{}) map[string]interface{} {

	out := make(map[string]interface{})

	// Loop over all our output structs
	for _, data := range interfaces {
		for key, value := range data {
			out[key] = value
		}
	}

	return out
}

func applyPrefix(data map[string]interface{}, prefix string) map[string]interface{} {

	// Create a new map
	newData := make(map[string]interface{})

	// Prepend our prefix onto the keys
	for key, value := range data {
		newData[prefix+key] = value
	}

	return newData
}
