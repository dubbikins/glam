package logging

import (
	"fmt"
)
var reset = "\033[0m"


func color( str, color string) string {
	return fmt.Sprintf("%s%s%s", color, str, reset)
}
func Red(str string) string {
	return color(str, "\033[1;31m")
}
func Green(str string) string {
	return color(str, "\033[1;32m")
}
func Yellow(str string) string {
	return color(str, "\033[1;33m")
}
func Blue(str string) string {
	return color(str, "\033[1;34m")
}
func Purple(str string) string {
	return color(str, "\033[1;35m")
}
func Cyan(str string) string {
	return color(str, "\033[1;36m")
}
func Gray(str string) string {
	return color(str, "\033[1;37m")
}
func Magenta(str string) string {
	return color(str, "\033[1;95m")
}
func White(str string) string {
	return color(str, "\033[1;97m")
}






