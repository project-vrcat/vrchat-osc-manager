package plugin

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/logger"

	"github.com/pelletier/go-toml"
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
		Dir         string `json:"-" toml:"-"` // plugin directory
	}
	Entrypoint struct {
		cmd        *exec.Cmd
		done       chan error
		Executable string   `json:"executable" toml:"executable"`
		Args       []string `json:"args" toml:"args"`
	}
)

func NewPlugin(dir, wsAddr string) (*Plugin, error) {
	p := &Plugin{Dir: dir, wsAddr: wsAddr}
	if err := p.load(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Plugin) load() error {
	t, err := toml.LoadFile(filepath.Join(p.Dir, "manifest.toml"))
	if err != nil {
		return err
	}
	if err = t.Unmarshal(p); err != nil {
		f, err := ioutil.ReadFile(filepath.Join(p.Dir, "manifest.json"))
		if err != nil {
			return errors.New("can not open manifest file: " + p.Dir)
		}
		if err = json.Unmarshal(f, p); err != nil {
			return err
		}
	}
	p.Enabled = false
	return nil
}

func (p *Plugin) Init() (err error) {
	if err = p.Entrypoint.checkExecutable(); err != nil {
		return err
	}
	if p.Install != nil {
		// check if the plugin is already installed
		if _, err = os.Stat(filepath.Join(p.Dir, ".installed")); err != nil {
			if err = p.Install.checkExecutable(); err != nil {
				return err
			}
			if err = p.Install.Start(p.Dir, p.Name, p.wsAddr); err != nil {
				return err
			}
			if err := p.Install.Wait(); err != nil {
				log.Println(p.Name, "installation failed:", err)
				return err
			}
			_ = ioutil.WriteFile(filepath.Join(p.Dir, ".installed"), nil, 0644)
		}
	}
	return
}

func (p *Plugin) Start() (err error) {
	if err = p.Entrypoint.Start(p.Dir, p.Name, p.wsAddr); err != nil {
		return err
	}
	return nil
}

func (p *Plugin) Stop() (err error) {
	if err = p.Entrypoint.Stop(); err != nil {
		return err
	}
	return nil
}

func (e *Entrypoint) checkExecutable() error {
	f := filepath.Join(config.C.RuntimeDir, e.Executable)
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

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	e.done = make(chan error)
	done := make(chan error)
	defer func() {
		go func() {
			done <- cmd.Wait()
		}()
	}()

	scanner := func(r io.ReadCloser, _log func(string, string)) {
		reader := bufio.NewReader(r)
		for {
			select {
			case err := <-done:
				e.done <- err
				return
			default:
				line, err := reader.ReadString('\n')
				if err == nil {
					_log(pluginName, line)
				}
			}
		}
	}

	go scanner(stdout, e.log)
	go scanner(stderr, e.error)

	return cmd.Start()
}

func (e *Entrypoint) Wait() *exec.ExitError {
	if e.cmd != nil && e.done != nil {
		if err := <-e.done; err != nil {
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

func (e *Entrypoint) error(pluginName, err string) {
	_err := strings.TrimSpace(err)
	logger.Plugin(pluginName).Error(_err)
}

func (e *Entrypoint) log(pluginName, line string) {
	logger.Plugin(pluginName).Print(line)
}
