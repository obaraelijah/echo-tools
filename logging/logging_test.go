package logging_test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestInitialize(t *testing.T) {
	filter := regexp.MustCompile("\u001B\\[\\d+m")
	fmt.Println(filter.MatchString("\u001B[31m"))
}
