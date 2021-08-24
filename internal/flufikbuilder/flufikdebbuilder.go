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

func (flufikDeb *FlufikDEBBuilder) metaData(flufikMeta flufikinfo.FlufikPackageMeta) flufikdeb.FlufikDebMetaData {
	return flufikdeb.FlufikDebMetaData{
		Package:      flufikMeta.Name,
		Version:      flufikMeta.Version,
		Maintainer:   flufikMeta.Maintainer,
		Summary:      flufikMeta.Summary,
		Description:  flufikMeta.Description,
		Architecture: flufikDeb.arch(),
		Homepage:     flufikMeta.URL,
	}
}

func (flufikDeb *FlufikDEBBuilder) DirToDebFile(flufikInfo flufikinfo.FlufikPackageDir) flufikdeb.FlufikDebFile {
	return flufikdeb.FlufikDebFile{
		Name:  flufikInfo.Destination,
		Mode:  flufikInfo.Mode + 040000,
		Owner: flufikInfo.Owner,
		Group: flufikInfo.Group,
		Type:  flufikdeb.Directory,
	}
}

func (flufikDeb *FlufikDEBBuilder) FileToDebFile(tName string, flufikInfo flufikinfo.FlufikPackageFile) (flufikdeb.FlufikDebFile, error) {
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
		flufikDeb.configurationFiles = append(flufikDeb.configurationFiles, flufikInfo.Destination)
	}

	return flufikdeb.FlufikDebFile{
		Name:  flufikInfo.Destination,
		Body:  flufikInfo.FileData(),
		Mode:  flufikInfo.FileMode(),
		Owner: flufikInfo.Owner,
		Group: flufikInfo.Group,
		MTime: flufikInfo.FileMTime(),
		Type:  fileType,
	}, nil

}

func (flufikDeb *FlufikDEBBuilder) arch() string {
	meta := flufikDeb.packageInfo.Meta
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

func (flufikDeb *FlufikDEBBuilder) FileName() (string, error) {
	meta := flufikDeb.packageInfo.Meta

	if meta.Name == "" {
		return "", fmt.Errorf("undefined package name")
	} else if meta.Version == "" {
		return "", fmt.Errorf("undefined package version")
	}
	return fmt.Sprintf("%s_%s_%s.deb", meta.Name, meta.Version, flufikDeb.arch()), nil
}

func (flufikDeb *FlufikDEBBuilder) Build(writer io.Writer) error {
	var (
		flufikDebPkg *flufikdeb.FlufikDeb
		err          error
	)
	if flufikDebPkg, err = flufikdeb.NewDeb(flufikDeb.metaData(flufikDeb.packageInfo.Meta)); err != nil {
		return err
	}

	for _, dir := range flufikDeb.packageInfo.Directory {
		flufikDebPkg.AddFile(flufikDeb.DirToDebFile(dir))
	}

	for tName, fList := range flufikDeb.packageInfo.Files {
		for _, file := range fList {
			if dFile, err := flufikDeb.FileToDebFile(tName, file); err == nil {
				flufikDebPkg.AddFile(dFile)
			} else {
				return err
			}
		}
	}

	flufikDebPkg.AddPreIn(flufikDeb.packageInfo.PreInScript())
	flufikDebPkg.AddPostIn(flufikDeb.packageInfo.PostInScript())
	flufikDebPkg.AddPreUn(flufikDeb.packageInfo.PreUnScript())
	flufikDebPkg.AddPostUn(flufikDeb.packageInfo.PostUnScript())

	//Signature part
	flufikDebPkg.AddSignatureKey(flufikDeb.packageInfo.AddSignatureKey())
	flufikDebPkg.AddSignatureType(flufikDeb.packageInfo.AddSignatureType())
	flufikDebPkg.AddSignaturePassPhrase(flufikDeb.packageInfo.AddSignaturePassPhrase())

	for _, dep := range flufikDeb.packageInfo.Dependencies {
		if err = flufikDebPkg.Depends.Set(dep.FlufikDEBFormat()); err != nil {
			return err
		}
	}
	return flufikDebPkg.Write(writer)
}

func NewFlufikDebBuilder(flkPkgInfo *flufikinfo.FlufikPackage) FlufikPackageBuilder {
	return &FlufikDEBBuilder{
		packageInfo: flkPkgInfo,
	}
}
