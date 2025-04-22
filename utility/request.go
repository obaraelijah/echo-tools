package utility

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
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
		return err
	}

	err = json.Unmarshal(b, form)
	if err != nil {
		return err
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		e = e.Elem()
	}

	var missing []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		required := field.Tag.Get(tagName) == "required"
		jsonName := field.Tag.Get("json")
		if s := strings.Split(jsonName, ","); len(s) > 1 {
			jsonName = s[0]
		}

		if required && e.Field(i).IsNil() {
			missing = append(missing, jsonName)
		}
	}
	if len(missing) == 1 {
		return fmt.Errorf("parameter %s is missing but required", missing[0])
	} else if len(missing) > 1 {
		return fmt.Errorf("parameter %s are missing but required", strings.Join(missing, ", "))
	}

	return nil
}
