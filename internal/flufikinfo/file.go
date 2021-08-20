package flufikinfo

import (
	"io/ioutil"
	"os"
	"time"
)

type FlufikPackageFile struct {
	Destination string    `yaml:"destination,omitempty"`
	Source      string    `yaml:"source"`
	Body        string    `yaml:"body"`
	Mode        uint      `yaml:"mode"`
	Owner       string    `yaml:"owner"`
	Group       string    `yaml:"group"`
	MTime       time.Time `yaml:"mtime"`
}

type FlufikPackageDir struct {
	Destination string `yaml:"destination,omitempty"`
	Mode        uint   `yaml:"mode"`
	Owner       string `yaml:"owner"`
	Group       string `yaml:"group"`
}

func (flufikPkgFile *FlufikPackageFile) FileData() []byte {
	if data, err := ioutil.ReadFile(flufikPkgFile.Source); err == nil {
		return data
	} else if flufikPkgFile.Body != "" {
		return []byte(flufikPkgFile.Body)
	} else {
		return make([]byte, 0)
	}
}

func (flufikPkgFile *FlufikPackageFile) FileMode() uint {
	if stat, err := os.Stat(flufikPkgFile.Source); err == nil && !stat.IsDir() {
		return uint(stat.Mode())
	} else {
		return flufikPkgFile.Mode
	}
}

func (flufikPkgFile *FlufikPackageFile) FileMTime() time.Time {
	if stat, err := os.Stat(flufikPkgFile.Source); err == nil && !stat.IsDir() {
		return stat.ModTime()
	} else {
		return flufikPkgFile.MTime
	}
}
