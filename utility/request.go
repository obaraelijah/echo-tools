package utility

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	funcgo "github.com/obaraelijah/funcgo"
)

var json = jsoniter.Config{
	EscapeHTML:    true,
	CaseSensitive: true,
}.Froze()

const tagName = "echotools"

// ValidateJsonForm use this method to validate a json request. Annotate your struct with `echotools:"required"` to
// mark the field as required.
func ValidateJsonForm(c echo.Context, form interface{}) error {
	t := reflect.TypeOf(form)
	e := reflect.ValueOf(form)

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return errors.New("error while reading body")
	}

	err = json.Unmarshal(b, form)
	if err != nil {
		return errors.New("error while decoding json")
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		e = e.Elem()
	}

	var missing []string
	var notEmptyViolated []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		fieldElem := e.Field(i)
		tags := strings.Split(field.Tag.Get(tagName), ";")
		cleaned := funcgo.Map(func(elem string) string { return strings.TrimSpace(elem) })(tags)
		required := funcgo.Any(func(elem string) bool { return elem == "required" })(cleaned)
		notEmpty := funcgo.Any(func(elem string) bool { return elem == "not empty" })(cleaned)
		jsonName := field.Tag.Get("json")
		if s := strings.Split(jsonName, ","); len(s) > 1 {
			jsonName = s[0]
		}

		isPointer := fieldType.Kind() == reflect.Ptr
		isNil := false

		// Required validation -> Can only be done on pointer
		if required {
			if isPointer {
				if fieldElem.IsNil() {
					isNil = true
					missing = append(missing, jsonName)
				}
			} else {
				// As this was probably not intended, output warnings
				c.Logger().Warnf("echotools required tag set on a non-pointer field: %s", jsonName)
			}
		}

		// Not empty validation -> Can only be done on string and *string
		// Test if not isNil -> dereference a nil pointer is evil (and panics)
		if !isNil && notEmpty {
			// Dereference pointer of needed
			if isPointer {
				fieldElem = fieldElem.Elem()
				fieldType = fieldElem.Type()
			}
			// Check if field type is string
			if fieldType.Kind() == reflect.String {
				if fieldElem.String() == "" {
					notEmptyViolated = append(notEmptyViolated, jsonName)
				}
			} else {
				// As this was probably not intended, output warnings
				c.Logger().Warnf("echotools not empty tag set on a non-string field: %s", jsonName)
			}
		}
	}
	if len(missing) == 1 {
		return fmt.Errorf("parameter %s is missing but required", missing[0])
	} else if len(missing) > 1 {
		return fmt.Errorf("parameter %s are missing but required", strings.Join(missing, ", "))
	}

	if len(notEmptyViolated) > 0 {
		name := strings.Join(notEmptyViolated, ", ")
		return fmt.Errorf("parameter %s must not be empty", name)
	}

	return nil
}
