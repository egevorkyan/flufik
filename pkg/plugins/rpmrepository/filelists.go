package rpmrepository

import "encoding/xml"

const filelistsXmlns = "http://linux.duke.edu/metadata/filelists"

type Filelists struct {
	XMLName  xml.Name           `xml:"filelists"`
	Packages int                `xml:"packages"`
	Package  []FilelistsPackage `xml:"package"`
	Xmlns    string             `xml:"xmlns,attr"`
}

type FilelistsPackage struct {
	File []File `xml:"file"`
	Package
}

func (r *RpmRepo) GetFilelists(packages PackageInfos) Filelists {
	filelists := Filelists{
		Packages: len(packages),
		Xmlns:    filelistsXmlns,
		Package:  []FilelistsPackage{},
	}

	for checksum, p := range packages {
		pkgversion := Version{
			Epoch:   p.Epoch(),
			Version: p.Version(),
			Release: p.Release(),
		}
		pkgdata := FilelistsPackage{
			Package: Package{
				Architecture: p.Architecture(),
				Pkgid:        checksum,
				Name:         p.Name(),
				Version:      pkgversion,
			},
			File: []File{},
		}
		for _, f := range p.Files() {
			file := File{
				Value: f,
			}
			//if f.IsDir() {
			//	file.Type = "dir"
			//}
			pkgdata.File = append(pkgdata.File, file)
		}
		filelists.Package = append(filelists.Package, pkgdata)
	}
	return filelists
}
