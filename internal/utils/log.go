package utils

import (
	"fmt"
	"strings"
)

func Color(color uint8, text ...string) string {
	return fmt.Sprintf("\033[0;0;%dm%s\033[0m", color, strings.Join(text, ""))
}
