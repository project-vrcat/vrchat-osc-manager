package plugin

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func NewPlugin(dir, wsAddr string) (*Plugin, error) {
	p := &Plugin{Dir: dir, wsAddr: wsAddr}
	if err := p.load(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Plugin) load() error {
	_, err := toml.DecodeFile(filepath.Join(p.Dir, "manifest.toml"), p)
	if err != nil {
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
	if err = p.Entrypoint.checkExecutable(p.Dir); err != nil {
		return err
	}
	if p.Install != nil {
		// check if the plugin is already installed
		if _, err = os.Stat(filepath.Join(p.Dir, ".installed")); err != nil {
			if err = p.Install.checkExecutable(p.Dir); err != nil {
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
