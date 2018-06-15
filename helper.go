package gotpl

import (
	"html/template"
)

// define build in helpers

func setFunc(data map[string]interface{}, key string, value interface{}) template.JS {
	data[key] = value
	return template.JS("")
}

func appendFunc(data map[string]interface{}, arr string, value interface{}) template.JS {
	if data[arr] == nil {
		data[arr] = []interface{}{value}
	} else {
		data[arr] = append(data[arr].([]interface{}), value)
	}
	return template.JS("")
}

func rawFunc(text string) template.HTML {
	return template.HTML(text)
}

// FuncMap return helpers
// set: set a variable in given context `{{ set . "key" "value" }}``
// append: append a variable to an array, or create one `{{ append . "arr" "value" }}``
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"set":    setFunc,
		"append": appendFunc,
		"raw":    rawFunc,
	}
}
