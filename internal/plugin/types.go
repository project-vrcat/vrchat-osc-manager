package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
)

type (
	Plugin struct {
		Name        string      `json:"name" toml:"name"`
		Author      string      `json:"author" toml:"author"`
		Version     string      `json:"version" toml:"version"`
		Description string      `json:"description" toml:"description"`
		HomePage    string      `json:"homepage" toml:"homepage"`
		Icon        string      `json:"icon" toml:"icon"`
		OptionsPage string      `json:"options_page" toml:"options_page"`
		Inputs      []string    `json:"inputs" toml:"inputs"`
		Entrypoint  Entrypoint  `json:"entrypoint" toml:"entrypoint"`
		Install     *Entrypoint `json:"install" toml:"install"`
		Enabled     bool        `json:"enabled" toml:"enabled"`
		wsAddr      string
		Dir         string // plugin directory
	}
	Entrypoint struct {
		cmd        *exec.Cmd
		Executable string   `json:"executable" toml:"executable"`
		Args       []string `json:"args" toml:"args"`
	}
)

func (e *Entrypoint) checkExecutable(dir string) error {
	f := filepath.Join(dir, e.Executable)
	_, err := os.Stat(f)
	if err != nil {
		_, err = exec.LookPath(e.Executable)
		return err
	}
	e.Executable, _ = filepath.Abs(f)
	return nil
}

func (e *Entrypoint) Start(pluginDir, pluginName, wsAddr string) error {
	cmd := exec.Command(e.Executable, e.Args...)
	cmd.Env = []string{
		"VRCOSCM_WS_ADDR=" + wsAddr,
		"VRCOSCM_PLUGIN=" + pluginName,
	}
	e.cmd = cmd
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func (e *Entrypoint) Wait() *exec.ExitError {
	if e.cmd != nil {
		if err := e.cmd.Wait(); err != nil {
			return err.(*exec.ExitError)
		}
	}
	return nil
}

func (e *Entrypoint) Stop() error {
	if e.cmd != nil {
		if e.cmd.Process != nil {
			return e.cmd.Process.Kill()
		}
	}
	return nil
}
