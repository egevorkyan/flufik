package flufikdeb

import (
	"archive/tar"
	"time"
)

type FileType int32

const (
	Directory FileType = 1 << iota >> 1
	GenericFile
	ConfigFile
	//ExcludeFile - future logic
)

type FlufikDebFile struct {
	Name  string
	Body  []byte
	Mode  uint
	Owner string
	Group string
	MTime time.Time
	Type  FileType
}

func (flufikDebF *FlufikDebFile) isDir() bool {
	return flufikDebF.Type == Directory
}

func (flufikDebF *FlufikDebFile) isConfig() bool {
	return (flufikDebF.Type & ConfigFile) == ConfigFile
}

func (flufikDebF *FlufikDebFile) tarTypeFlag() byte {
	if flufikDebF.isDir() {
		return tar.TypeDir
	} else {
		return tar.TypeReg
	}
}
