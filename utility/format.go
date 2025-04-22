package utility

import (
	"fmt"
)

func PPrintln(a interface{}) {
	indent, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println(string(indent))
}
