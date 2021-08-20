package flufikdeb

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Relations []string

type FlufikDebMetaData struct {
	Package,
	Version,
	Maintainer,
	Summary,
	Description,
	Section,
	Priority string
	Essential bool
	Architecture,
	Origin,
	Bugs,
	Homepage string
	Tag    []string
	Source string
	Depends,
	PreDepends,
	Recommends,
	Suggests,
	Breaks,
	Conflicts,
	Replaces,
	Provides Relations
}

const metaHeader = `Package: %s
Version: %s
Architecture: %s
Maintainer: %s
Homepage: %s
Depends: %s
Description: %s
%s`

func (rel *Relations) Set(pkg string) error {
	for _, relation := range *rel {
		if relation == pkg {
			return nil
		}
	}

	*rel = append(*rel, pkg)
	return nil
}

func (flufikMetaData *FlufikDebMetaData) depends() string {
	return strings.Join(flufikMetaData.Depends, ", ")
}

func (flufikMetaData *FlufikDebMetaData) descForm() string {
	return flufikMetaData.Description
}

func (flufikMetaData *FlufikDebMetaData) description() string {
	out := bytes.NewBufferString("")
	in := bufio.NewReader(bytes.NewBufferString(flufikMetaData.Description))

	for {
		if line, _, err := in.ReadLine(); err != io.EOF {
			_, _ = fmt.Fprintf(out, " %s\n", string(line))
		} else {
			break
		}
	}
	return out.String()
}

func (flufikMetaData *FlufikDebMetaData) MakeControl() []byte {
	return []byte(fmt.Sprintf(
		metaHeader,
		flufikMetaData.Package,
		flufikMetaData.Version,
		flufikMetaData.Architecture,
		flufikMetaData.Maintainer,
		flufikMetaData.Homepage,
		flufikMetaData.depends(),
		flufikMetaData.Summary,
		flufikMetaData.description(),
	))
}
