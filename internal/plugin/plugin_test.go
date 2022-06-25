package plugin

import (
	"testing"
	"time"
)

func TestPlugin_Init(t *testing.T) {
	p, err := NewPlugin("../../plugins/example-plugin", "localhost:8787")
	if err != nil {
		t.Error(err)
		return
	}
	if err = p.Init(); err != nil {
		t.Fatal(err)
	}
}

func TestPlugin_Start(t *testing.T) {
	p, err := NewPlugin("../../plugins/example-plugin", "localhost:8787")
	if err != nil {
		t.Error(err)
		return
	}
	if err = p.Init(); err != nil {
		t.Fatal(err)
		return
	}
	if err = p.Start(); err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(time.Second * 3)
	if err = p.Stop(); err != nil {
		t.Fatal(err)
		return
	}
}
