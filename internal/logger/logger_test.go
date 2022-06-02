package logger

import "testing"

func TestLogger_Error(t *testing.T) {
	App("App").Error("Boom!")
	Plugin("Plugin").Error("Boom!")
}

func TestLogger_Print(t *testing.T) {
	App("App").Print("Hello", "\n")
	Plugin("Plugin").Print("Hello", "\n")
}

func TestLogger_Println(t *testing.T) {
	App("App").Println("Hello")
	Plugin("Plugin").Println("Hello")
}

func TestLogger_Printf(t *testing.T) {
	App("App").Printf("Hello %s\n", "World")
	Plugin("Plugin").Printf("Hello %s\n", "World")
}
