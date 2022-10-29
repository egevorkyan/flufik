package flufikbuilder

import (
	"fmt"
	"github.com/egevorkyan/flufik/internal/flufikinfo"
	"github.com/egevorkyan/flufik/pkg/flufikdeb"
	"io"
)

type FlufikDEBBuilder struct {
	FlufikPackageBuilder
	packageInfo        *flufikinfo.FlufikPackage
	configurationFiles []string
}

func (d *FlufikDEBBuilder) metaData(flufikMeta flufikinfo.FlufikPackageMeta) flufikdeb.FlufikDebMetaData {
	return flufikdeb.FlufikDebMetaData{
		Package:      flufikMeta.Name,
		Version:      flufikMeta.Version,
		Release:      flufikMeta.Release,
		Maintainer:   flufikMeta.Maintainer,
		Summary:      flufikMeta.Summary,
		Description:  flufikMeta.Description,
		Architecture: d.arch(),
		Homepage:     flufikMeta.URL,
	}
}

func (d *FlufikDEBBuilder) DirToDebFile(flufikInfo flufikinfo.FlufikPackageDir) flufikdeb.FlufikDebFile {
	return flufikdeb.FlufikDebFile{
		Name:  flufikInfo.Destination,
		Mode:  flufikInfo.Mode + 040000,
		Owner: flufikInfo.Owner,
		Group: flufikInfo.Group,
		Type:  flufikdeb.Directory,
	}
}

func (d *FlufikDEBBuilder) FileToDebFile(tName string, flufikInfo flufikinfo.FlufikPackageFile) (flufikdeb.FlufikDebFile, error) {
	fileType := flufikdeb.GenericFile

	switch tName {
	case "generic":
		fileType = flufikdeb.GenericFile
	case "config":
		fileType = flufikdeb.ConfigFile
	case "doc":
		fileType = flufikdeb.GenericFile
	case "not_use":
		fileType = flufikdeb.GenericFile
	case "missing_ok":
		fileType = flufikdeb.GenericFile
	case "no_replace":
		fileType = flufikdeb.ConfigFile
	case "spec":
		fileType = flufikdeb.GenericFile
	case "ghost":
		fileType = flufikdeb.GenericFile
	case "license":
		fileType = flufikdeb.GenericFile
	case "readme":
		fileType = flufikdeb.GenericFile
	case "exclude":
		fileType = flufikdeb.GenericFile
	default:
		return flufikdeb.FlufikDebFile{}, fmt.Errorf("unexpected file type: %s", tName)

	}

	if (fileType & flufikdeb.ConfigFile) != 0 {
		d.configurationFiles = append(d.configurationFiles, flufikInfo.Destination)
	}

	body, err := flufikInfo.FileData()
	if err != nil {
		return flufikdeb.FlufikDebFile{}, err
	}

	return flufikdeb.FlufikDebFile{
		Name:  flufikInfo.Destination,
		Body:  body,
		Mode:  flufikInfo.FileMode(),
		Owner: flufikInfo.Owner,
		Group: flufikInfo.Group,
		MTime: flufikInfo.FileMTime(),
		Type:  fileType,
	}, nil

}

func (d *FlufikDEBBuilder) arch() string {
	meta := d.packageInfo.Meta
	arch := "all"
	if meta.Arch != "" {
		switch meta.Arch {
		case "386":
			arch = "i386"
		case "x86_64":
			arch = "amd64"
		case "noarch":
			arch = "all"
		default:
			arch = meta.Arch
		}
	}
	return arch
}

func (d *FlufikDEBBuilder) FileName() (string, error) {
	meta := d.packageInfo.Meta

	if meta.Name == "" {
		return "", fmt.Errorf("undefined package name")
	} else if meta.Version == "" {
		return "", fmt.Errorf("undefined package version")
	} else if meta.Release == "" {
		return "", fmt.Errorf("undefined package release")
	}
	return fmt.Sprintf("%s_%s-%s_%s.deb", meta.Name, meta.Version, meta.Release, d.arch()), nil
}

func (d *FlufikDEBBuilder) Build(writer io.Writer) error {
	var (
		fluffDebPkg *flufikdeb.FlufikDeb
		err         error
	)
	if fluffDebPkg, err = flufikdeb.NewDeb(d.metaData(d.packageInfo.Meta)); err != nil {
		return err
	}

	for _, dir := range d.packageInfo.Directory {
		fluffDebPkg.AddFile(d.DirToDebFile(dir))
	}

	for tName, fList := range d.packageInfo.Files {
		for _, file := range fList {
			if dFile, err := d.FileToDebFile(tName, file); err == nil {
				fluffDebPkg.AddFile(dFile)
			} else {
				return err
			}
		}
	}

	fluffDebPkg.AddPreIn(d.packageInfo.PreInScript())
	fluffDebPkg.AddPostIn(d.packageInfo.PostInScript())
	fluffDebPkg.AddPreUn(d.packageInfo.PreUnScript())
	fluffDebPkg.AddPostUn(d.packageInfo.PostUnScript())

	//Signature part
	fluffDebPkg.AddSignatureKey(d.packageInfo.AddSignatureKey())
	fluffDebPkg.AddSignatureType(d.packageInfo.AddSignatureType())

	for _, dep := range d.packageInfo.Dependencies {
		if err = fluffDebPkg.Depends.Set(dep.FlufikDEBFormat()); err != nil {
			return err
		}
	}
	return fluffDebPkg.Write(writer)
}

func NewFlufikDebBuilder(flkPkgInfo *flufikinfo.FlufikPackage) FlufikPackageBuilder {
	return &FlufikDEBBuilder{
		packageInfo: flkPkgInfo,
	}
}
