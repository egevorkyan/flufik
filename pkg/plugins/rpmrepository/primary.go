package rpmrepository

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
)

const primaryXmlns = "http://linux.duke.edu/metadata/common"
const primaryXmlnsRpm = "http://linux.duke.edu/metadata/rpm"

type Primary struct {
	XMLName  xml.Name         `xml:"metadata"`
	Packages int              `xml:"packages,attr"`
	Package  []PrimaryPackage `xml:"package"`
	Xmlns    string           `xml:"xmlns,attr"`
	Xmlnsrpm string           `xml:"xmlns:rpm,attr"`
}

type PrimaryPackage struct {
	Type         string        `xml:"type,attr"`
	Name         string        `xml:"name"`
	Architecture string        `xml:"arch"`
	Version      Version       `xml:"version"`
	Checksum     Checksum      `xml:"checksum"`
	Summary      string        `xml:"summary"`
	Description  string        `xml:"description"`
	Packager     string        `xml:"packager"`
	URL          string        `xml:"url"`
	Time         PrimaryTime   `xml:"time"`
	Size         PrimarySize   `xml:"size"`
	Format       PrimaryFormat `xml:"format"`
	Location     Location      `xml:"location"`
}

type PrimaryTime struct {
	File  int64 `xml:"file,attr"`
	Build int64 `xml:"build,attr"`
}

type PrimarySize struct {
	Package   uint64 `xml:"package,attr"`
	Installed uint64 `xml:"installed,attr"`
	Archived  uint64 `xml:"archived,attr"`
}

type PrimaryFormat struct {
	License     string                   `xml:"rpm:license"`
	Vendor      string                   `xml:"rpm:vendor"`
	Groups      []string                 `xml:"rpm:group"`
	Buildhost   string                   `xml:"rpm:buildhost"`
	SourceRPM   string                   `xml:"rpm:sourcerpm"`
	HeaderRange PrimaryFormatHeaderRange `xml:"rpm:heander-range"`
	Provides    []PrimaryFormatEntry     `xml:"rpm:provides>rpm:entry"`
	Requires    []PrimaryFormatEntry     `xml:"rpm:requires>rpm:entry"`
	Conflicts   []PrimaryFormatEntry     `xml:"rpm:conflicts>rpm:entry"`
	Files       []File                   `xml:"file"`
}

type PrimaryFormatHeaderRange struct {
	Start uint64 `xml:"start,attr"`
	End   uint64 `xml:"end,attr"`
}

type PrimaryFormatEntry struct {
	Name    string `xml:"name,attr"`
	Flags   string `xml:"flags,attr,omitempty"`
	Epoch   string `xml:"epoch,attr,omitempty"`
	Version string `xml:"ver,attr,omitempty"`
	Release string `xml:"rel,attr,omitempty"`
	Pre     string `xml:"pre,attr,omitempty"`
}

func (r *RpmRepo) GetPrimary(packages PackageInfos) Primary {
	fileRegex := []*regexp.Regexp{
		regexp.MustCompile(".*bin/.*"),
		regexp.MustCompile("^/etc/.*"),
		regexp.MustCompile("^/usr/lib/sendmail$"),
	}
	primary := Primary{
		Packages: len(packages),
		Xmlns:    primaryXmlns,
		Xmlnsrpm: primaryXmlnsRpm,
		Package:  []PrimaryPackage{},
	}

	for checksum, p := range packages {
		pkgversion := Version{
			Epoch:   p.Epoch(),
			Version: p.Version(),
			Release: p.Release(),
		}
		pkgsum := Checksum{
			Value: checksum,
			Type:  "sha256",
			Pkgid: "YES",
		}
		pkgtime := PrimaryTime{
			File:  p.FileTime().Unix(),
			Build: p.BuildTime().Unix(),
		}
		// TODO: Sizes seem not to work quite well
		pkgsize := PrimarySize{
			Package:   p.FileSize(),
			Installed: p.Size(),
			Archived:  p.ArchiveSize(),
		}
		pkgformatheaderrange := PrimaryFormatHeaderRange{
			Start: p.HeaderStart(),
			End:   p.HeaderEnd(),
		}
		pkgformat := PrimaryFormat{
			License:     p.License(),
			Vendor:      p.Vendor(),
			Groups:      p.Groups(),
			Buildhost:   p.BuildHost(),
			SourceRPM:   p.SourceRPM(),
			HeaderRange: pkgformatheaderrange,
			Provides:    []PrimaryFormatEntry{},
			Requires:    []PrimaryFormatEntry{},
			Conflicts:   []PrimaryFormatEntry{},
			Files:       []File{},
		}
		for _, p := range p.Provides() {
			provided := PrimaryFormatEntry{
				Name:    p.Name(),
				Epoch:   strconv.Itoa(p.Epoch()),
				Release: p.Release(),
				Version: p.Version(),
				Flags:   r.ReadFlags(p.Flags()),
			}
			pkgformat.Provides = append(pkgformat.Provides, provided)
		}
		for _, r := range p.Requires() {
			requirement := PrimaryFormatEntry{
				Name: r.Name(),
			}
			pkgformat.Requires = append(pkgformat.Requires, requirement)
		}
		for _, c := range p.Conflicts() {
			confilct := PrimaryFormatEntry{
				Name: c.Name(),
			}
			pkgformat.Conflicts = append(pkgformat.Conflicts, confilct)
		}
		for _, f := range p.Files() {
			file := File{
				Value: f,
			}
			//if f.IsDir() {
			//	// TODO: The if does not quite work
			//	file.Type = "dir"
			//}
			for _, re := range fileRegex {
				if re.MatchString(f) {
					pkgformat.Files = append(pkgformat.Files, file)
					break
				}
			}
		}
		pkg := PrimaryPackage{
			Type:         "rpm",
			Name:         p.Name(),
			Architecture: p.Architecture(),
			Version:      pkgversion,
			Checksum:     pkgsum,
			Summary:      p.Summary(),
			Description:  p.Description(),
			Packager:     p.Packager(),
			URL:          p.URL(),
			Time:         pkgtime,
			Size:         pkgsize,
			Format:       pkgformat,
			Location: Location{
				Href: fmt.Sprintf("%s.rpm", p.Path),
			},
		}
		primary.Package = append(primary.Package, pkg)
	}
	return primary
}
