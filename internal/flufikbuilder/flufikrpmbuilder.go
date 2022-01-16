package flufikbuilder

import (
	"errors"
	"fmt"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/internal/flufikinfo"
	"github.com/google/rpmpack"
	"io"
)

type FlufikRPMBuilder struct {
	FlufikPackageBuilder
	packageInfo *flufikinfo.FlufikPackage
}

func (flufikRpm *FlufikRPMBuilder) rpmMetadata(flufikMeta flufikinfo.FlufikPackageMeta) rpmpack.RPMMetaData {
	return rpmpack.RPMMetaData{
		Name:        flufikMeta.Name,
		Summary:     flufikMeta.Summary,
		Description: flufikMeta.Description,
		Version:     flufikMeta.Version,
		Release:     flufikMeta.Release,
		Arch:        flufikMeta.Arch,
		OS:          flufikMeta.OS,
		Vendor:      flufikMeta.Vendor,
		URL:         flufikMeta.URL,
		Packager:    flufikMeta.Maintainer,
		Group:       "",
		Licence:     flufikMeta.License,
		BuildHost:   "",
		Compressor:  "",
		Epoch:       0,
		BuildTime:   flufikMeta.BuildTime(),
		Provides:    nil,
		Obsoletes:   nil,
		Suggests:    nil,
		Recommends:  nil,
		Requires:    nil,
		Conflicts:   nil,
	}
}

func (flufikRpm *FlufikRPMBuilder) dirToRPMFile(flufikInfo flufikinfo.FlufikPackageDir) rpmpack.RPMFile {
	return rpmpack.RPMFile{
		Name:  flufikInfo.Destination,
		Mode:  flufikInfo.Mode + 040000,
		Owner: flufikInfo.Owner,
		Group: flufikInfo.Group,
	}
}

func (flufikRpm *FlufikRPMBuilder) fileToRpmFile(tName string, flufikInfo flufikinfo.FlufikPackageFile) (rpmpack.RPMFile, error) {
	fileType := rpmpack.GenericFile

	switch tName {
	case "generic":
		fileType = rpmpack.GenericFile
	case "config":
		fileType = rpmpack.ConfigFile | rpmpack.NoReplaceFile
	case "doc":
		fileType = rpmpack.DocFile
	case "not_use":
		fileType = rpmpack.DoNotUseFile
	case "missing_ok":
		fileType = rpmpack.MissingOkFile
	case "no_replace":
		fileType = rpmpack.NoReplaceFile
	case "spec":
		fileType = rpmpack.SpecFile
	case "ghost":
		fileType = rpmpack.GhostFile
	case "license":
		fileType = rpmpack.LicenceFile
	case "readme":
		fileType = rpmpack.ReadmeFile
	case "exclude":
		fileType = rpmpack.ExcludeFile
	default:
		return rpmpack.RPMFile{}, errors.New("unexpected file type: " + tName)
	}

	return rpmpack.RPMFile{
		Name:  flufikInfo.Destination,
		Body:  flufikInfo.FileData(),
		Mode:  flufikInfo.FileMode(),
		Owner: flufikInfo.Owner,
		Group: flufikInfo.Group,
		MTime: uint32(flufikInfo.FileMTime().Unix()),
		Type:  fileType,
	}, nil
}

func (flufikRpm *FlufikRPMBuilder) FileName() (string, error) {
	flufikMeta := flufikRpm.packageInfo.Meta

	if flufikMeta.Name == "" {
		return "", errors.New("undefined package name")
	} else if flufikMeta.Version == "" {
		return "", errors.New("undefined package version")
	}

	release := ""
	if flufikMeta.Release != "" {
		release = "-" + flufikMeta.Release
	}

	arch := "noarch"
	if flufikMeta.Arch != "" {
		arch = flufikMeta.Arch
	}

	return fmt.Sprintf("%s-%s%s.%s.rpm", flufikMeta.Name, flufikMeta.Version, release, arch), nil
}

func (flufikRpm *FlufikRPMBuilder) Build(writer io.Writer) error {
	var (
		flufikRpmPkg *rpmpack.RPM
		err          error
	)
	if flufikRpmPkg, err = rpmpack.NewRPM(flufikRpm.rpmMetadata(flufikRpm.packageInfo.Meta)); err != nil {
		return err
	}

	for _, dir := range flufikRpm.packageInfo.Directory {
		flufikRpmPkg.AddFile(flufikRpm.dirToRPMFile(dir))
	}

	for tName, fList := range flufikRpm.packageInfo.Files {
		for _, file := range fList {
			if rpmFile, err := flufikRpm.fileToRpmFile(tName, file); err == nil {
				flufikRpmPkg.AddFile(rpmFile)
			} else {
				return err
			}
		}
	}

	flufikRpmPkg.AddPrein(flufikRpm.packageInfo.PreInScript())
	flufikRpmPkg.AddPostin(flufikRpm.packageInfo.PostInScript())
	flufikRpmPkg.AddPreun(flufikRpm.packageInfo.PreUnScript())
	flufikRpmPkg.AddPostun(flufikRpm.packageInfo.PostUnScript())

	for _, dep := range flufikRpm.packageInfo.Dependencies {
		if err = flufikRpmPkg.Requires.Set(dep.FlufikRPMFormat()); err != nil {
			return err
		}
	}
	if flufikRpm.packageInfo.Signature.PrivateKey != "" {
		flufikRpmPkg.SetPGPSigner(crypto.FlufikRpmSigner(flufikRpm.packageInfo.Signature.PrivateKey))

	}
	return flufikRpmPkg.Write(writer)
}

func NewFlufikRpmBuilder(flufikPkgInfo *flufikinfo.FlufikPackage) FlufikPackageBuilder {
	return &FlufikRPMBuilder{
		packageInfo: flufikPkgInfo,
	}
}
