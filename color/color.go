package color

import "fmt"

type Color string

const (
	reset  Color = "\033[0m"
	RED          = "\033[31m"
	GREEN        = "\033[32m"
	YELLOW       = "\033[33m"
	BLUE         = "\033[34m"
	PURPLE       = "\033[35m"
	CYAN         = "\033[36m"
	WHITE        = "\033[97m"
)

func Colorize(color Color, s string) string {
	return string(color) + s + string(reset)
}

func Print(color Color, s string) {
	fmt.Print(string(color) + s + string(reset))
}

func Printf(color Color, format string, args ...interface{}) {
	fmt.Print(string(color) + fmt.Sprintf(format, args...) + string(reset))
}

func Println(color Color, s string) {
	fmt.Println(string(color) + s + string(reset))
}
