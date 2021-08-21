package flufikbuilder

import "io"

type FlufikPackageBuilder interface {
	FileName() (string, error)
	Build(rPath io.Writer) error
}
