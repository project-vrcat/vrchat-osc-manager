package plugin

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func NewPlugin(dir string) (*Plugin, error) {
	p := &Plugin{dir: dir}
	if err := p.load(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Plugin) load() error {
	_, err := toml.DecodeFile(filepath.Join(p.dir, "manifest.toml"), p)
	if err != nil {
		f, err := ioutil.ReadFile(filepath.Join(p.dir, "manifest.json"))
		if err != nil {
			return errors.New("can not open manifest file: " + p.dir)
		}
		if err = json.Unmarshal(f, p); err != nil {
			return err
		}
	}
	p.Entrypoint.dir = p.dir
	p.Entrypoint.name = p.Name
	if p.Install != nil {
		p.Install.dir = p.dir
		p.Install.name = p.Name
	}
	return nil
}

func (p *Plugin) Init() (err error) {
	if err = p.Entrypoint.checkExecutable(); err != nil {
		return err
	}
	if p.Install != nil {
		// check if the plugin is already installed
		if _, err = os.Stat(filepath.Join(p.Install.dir, ".installed")); err != nil {
			if err = p.Install.checkExecutable(); err != nil {
				return err
			}
			if err = p.Install.Start(); err != nil {
				return err
			}
			if err := p.Install.Wait(); err != nil {
				log.Println("[plugin]", p.Name, "installation failed:", err)
				return err
			}
			_ = ioutil.WriteFile(filepath.Join(p.dir, ".installed"), nil, 0644)
		}
	}
	return
}

func (p *Plugin) Start() (err error) {
	if err = p.Entrypoint.Start(); err != nil {
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
