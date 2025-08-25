package lib

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateInsertQueryTagJson(name string, data interface{}) string {
	query := "INSERT INTO " + name + " ("
	values := "VALUES ("

	// Get the type and value of the input struct
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// Check if data is a pointer, if yes, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check if the field is a pointer
		if value.Kind() == reflect.Ptr {
			// If it's a pointer and it's nil, skip it
			if value.IsNil() {
				continue
			}
			// Dereference the pointer to get the actual value
			value = value.Elem()
		}

		// Get the field name and value
		fieldName := field.Tag.Get("json")
		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}
		fieldValue := fmt.Sprintf("%v", value.Interface())

		// Check if the field has omitempty tag and if the value is empty
		if omitempty := field.Tag.Get("omitempty"); omitempty != "" && fieldValue == "" {
			continue // Skip fields with empty values if they have omitempty tag
		}

		// Append the field name and value to the query
		query += fieldName + ", "
		values += "'" + fieldValue + "', "
	}

	// Remove the trailing comma and space
	query = strings.TrimSuffix(query, ", ")
	values = strings.TrimSuffix(values, ", ")

	// Complete the query
	query += ") " + values + ")"

	return query
}

func GenerateInsertQueryTagDb(name string, data interface{}) string {
	query := "INSERT INTO " + name + " ("
	values := "VALUES ("

	// Get the type and value of the input struct
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// Check if data is a pointer, if yes, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check if the field is a pointer
		if value.Kind() == reflect.Ptr {
			// If it's a pointer and it's nil, skip it
			if value.IsNil() {
				continue
			}
			// Dereference the pointer to get the actual value
			value = value.Elem()
		}

		// Get the field name and value
		fieldName := field.Tag.Get("db")
		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}
		fieldValue := fmt.Sprintf("%v", value.Interface())

		// Check if the field has omitempty tag and if the value is empty
		if omitempty := field.Tag.Get("omitempty"); omitempty != "" && fieldValue == "" {
			continue // Skip fields with empty values if they have omitempty tag
		}

		// Append the field name and value to the query
		query += fieldName + ", "
		values += "'" + fieldValue + "', "
	}

	// Remove the trailing comma and space
	query = strings.TrimSuffix(query, ", ")
	values = strings.TrimSuffix(values, ", ")

	// Complete the query
	query += ") " + values + ")"

	return query
}

func GenerateUpdateQueryTagJson(name string, data interface{}, id string) string {
	query := "UPDATE " + name + " SET "

	// Get the type and value of the input struct
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// Check if data is a pointer, if yes, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check if the field is a pointer
		if value.Kind() == reflect.Ptr {
			// If it's a pointer and it's nil, skip it
			if value.IsNil() {
				continue
			}
			// Dereference the pointer to get the actual value
			value = value.Elem()
		}

		// Get the field name and value
		fieldName := field.Tag.Get("json")
		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}
		fieldValue := fmt.Sprintf("%v", value.Interface())

		// Check if the field has omitempty tag and if the value is empty
		if omitempty := field.Tag.Get("omitempty"); omitempty != "" && fieldValue == "" {
			continue // Skip fields with empty values if they have omitempty tag
		}

		// Append the field name and value to the query
		query += fieldName + " = '" + fieldValue + "', "
	}

	// Remove the trailing comma and space
	query = strings.TrimSuffix(query, ", ")

	// Complete the query
	query += " WHERE id = " + id

	return query
}

func GenerateUpdateQueryTagDb(name string, data interface{}, id string) string {
	query := "UPDATE " + name + " SET "

	// Get the type and value of the input struct
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// Check if data is a pointer, if yes, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check if the field is a pointer
		if value.Kind() == reflect.Ptr {
			// If it's a pointer and it's nil, skip it
			if value.IsNil() {
				continue
			}
			// Dereference the pointer to get the actual value
			value = value.Elem()
		}

		// Get the field name and value
		fieldName := field.Tag.Get("db")
		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}
		fieldValue := fmt.Sprintf("%v", value.Interface())

		// Check if the field has omitempty tag and if the value is empty
		if omitempty := field.Tag.Get("omitempty"); omitempty != "" && fieldValue == "" {
			continue // Skip fields with empty values if they have omitempty tag
		}

		// Append the field name and value to the query
		query += fieldName + " = '" + fieldValue + "', "
	}

	// Remove the trailing comma and space
	query = strings.TrimSuffix(query, ", ")

	// Complete the query
	query += " WHERE id = " + id

	return query
}

func GenerateDeleteQuery(name string, id string) string {
	query := "DELETE FROM " + name + " WHERE id = " + id
	return query
}
