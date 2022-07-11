package rpmrepository

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"github.com/jfrog/go-rpm"
)

type Package struct {
	Version      Version `xml:"version"`
	Architecture string  `xml:"arch,attr"`
	Pkgid        string  `xml:"pkgid,attr"`
	Name         string  `xml:"name,attr"`
}

type Version struct {
	Epoch   int    `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Release string `xml:"rel,attr"`
}

type Checksum struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
	Pkgid string `xml:"pkgid,attr,omitempty"`
}

type Location struct {
	Href string `xml:"href,attr"`
}

type File struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

func (r *RpmRepo) ReadFlags(f int) string {
	var s string
	switch {
	case rpm.DepFlagLesserOrEqual == (f & rpm.DepFlagLesserOrEqual):
		s = "LE"

	case rpm.DepFlagLesser == (f & rpm.DepFlagLesser):
		s = "LT"

	case rpm.DepFlagGreaterOrEqual == (f & rpm.DepFlagGreaterOrEqual):
		s = "GE"

	case rpm.DepFlagGreater == (f & rpm.DepFlagGreater):
		s = "GT"

	case rpm.DepFlagEqual == (f & rpm.DepFlagEqual):
		s = "EQ"
	}
	return s
}

func (r *RpmRepo) GetXML(in interface{}) (out []byte, size int, checksum string, err error) {
	out, err = xml.MarshalIndent(in, "", "    ")
	if err != nil {
		return
	}
	out = []byte(xml.Header + string(out))

	sum := sha256.Sum256(out)
	checksum = fmt.Sprintf("%x", sum)

	size = len(out)
	return
}

func (r *RpmRepo) GetZip(in []byte) (out []byte, size int, checksum string) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(in)
	if err != nil {
		return nil, 0, ""
	}
	err = writer.Close()
	if err != nil {
		return nil, 0, ""
	}
	out = buf.Bytes()

	sum := sha256.Sum256(out)
	checksum = fmt.Sprintf("%x", sum)

	size = len(out)

	return
}
