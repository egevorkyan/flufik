package flufikinfo

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
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

func (flufikPkgFile *FlufikPackageFile) FileData() ([]byte, error) {
	absPath, err := core.FlufikMakePathAbs(flufikPkgFile.Source)
	if err != nil {
		return nil, err
	}
	if data, err := os.ReadFile(absPath); err == nil {
		return data, nil
	} else if flufikPkgFile.Body != "" {
		return []byte(flufikPkgFile.Body), nil
	} else {
		curDir, _ := os.Getwd()
		return nil, fmt.Errorf("path is wrong or file/directory does not exists, tried to reach from this workdir %s to target %s: %v", curDir, flufikPkgFile.Source, err)
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
