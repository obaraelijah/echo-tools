package utility

import (
	"encoding/json"
	"fmt"
)

func PPrintln(a interface{}) {
	indent, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println(string(indent))
}
