package logger

import (
	"fmt"
	"os"

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

func (l *logger) Error(err any) *logger {
	_err := color.BgRed.Render(color.FgWhite.Render(" ERROR "))
	fmt.Println(l.name+_err, err)
	return l
}

func (l *logger) Print(a ...any) *logger {
	v := fmt.Sprintln(append([]any{l.name}, a...)...)
	fmt.Print(v[:len(v)-1])
	return l
}

func (l *logger) Println(a ...any) *logger {
	fmt.Println(append([]any{l.name}, a...)...)
	return l
}

func (l *logger) Printf(format string, a ...any) *logger {
	fmt.Print(l.name, " ", fmt.Sprintf(format, a...))
	return l
}

func (l *logger) Fatal() {
	os.Exit(1)
}
