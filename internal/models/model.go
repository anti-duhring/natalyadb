package models

import (
	"reflect"
	"sort"
	"strings"
)

func ReadTags(s interface{}) {
	t := reflect.TypeOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("natdb")
		if tag != "" {
			println(tag)
		}
	}
}

func GetTypes(s interface{}) map[string]string {
	t := reflect.TypeOf(s)
	types := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := strings.ToLower(field.Name)
		types[name] = field.Type.Name()
	}

	return types
}

func GetSchema(s interface{}) []string {
	t := reflect.TypeOf(s)
	schema := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := strings.ToLower(field.Name)
		fieldType := field.Type.Name()
		schema = append(schema, name+":"+fieldType)
	}

	sort.Strings(schema)

	return schema
}
