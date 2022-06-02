package logger

import (
	"fmt"
	"github.com/gookit/color"
)

type (
	logger struct {
		name string
	}
)

func _new(name string, c color.Color) *logger {
	return &logger{name: c.Render("[" + name + "]")}
}

func App(name string) *logger {
	return _new(name, color.FgLightBlue)
}

func Plugin(name string) *logger {
	return _new(name, color.FgCyan)
}

func (l *logger) Error(err any) {
	_err := color.BgRed.Render(color.FgWhite.Render(" ERROR "))
	fmt.Println(_err+l.name, err)
}

func (l *logger) Print(a ...any) {
	v := fmt.Sprintln(append([]any{l.name}, a...)...)
	fmt.Print(v[:len(v)-1])
}

func (l *logger) Println(a ...any) {
	fmt.Println(append([]any{l.name}, a...)...)
}

func (l *logger) Printf(format string, a ...any) {
	fmt.Print(l.name, " ", fmt.Sprintf(format, a...))
}
