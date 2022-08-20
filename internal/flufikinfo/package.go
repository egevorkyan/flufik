package flufikinfo

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"time"
)

type FlufikPackage struct {
	Meta          FlufikPackageMeta              `yaml:"meta"`
	Directory     []FlufikPackageDir             `yaml:"directory"`
	Files         map[string][]FlufikPackageFile `yaml:"files"`
	PreInstall    []string                       `yaml:"preinstall"`
	PostInstall   []string                       `yaml:"postinstall"`
	PreUninstall  []string                       `yaml:"preuninstall"`
	PostUninstall []string                       `yaml:"postuninstall"`
	Dependencies  []FlufikDependency             `yaml:"dependencies"`
	Signature     FlufikDebSignature             `yaml:"signature"`
	sourceHome    string
}

type FlufikDebSignature struct {
	FlufikPackageSignature `yaml:",inline"`
	Type                   string `yaml:"type,omitempty" default:"origin"`
}

type FlufikPackageSignature struct {
	//PrivateKey string `yaml:"private_key,omitempty"`
	//PassPhrase string `yaml:"pass_phrase,omitempty"`
	PgpName string `yaml:"pgp_name,omitempty"`
}

func (flufikPkg *FlufikPackage) AppendPreIn(script string) {
	flufikPkg.PreInstall = append(flufikPkg.PreInstall, script+";")
}

func (flufikPkg *FlufikPackage) AppendPostIn(script string) {
	flufikPkg.PostInstall = append(flufikPkg.PostInstall, script+";")
}

func (flufikPkg *FlufikPackage) AppendPreUn(script string) {
	flufikPkg.PreUninstall = append(flufikPkg.PreUninstall, script+";")
}

func (flufikPkg *FlufikPackage) AppendPostUn(script string) {
	flufikPkg.PostUninstall = append(flufikPkg.PostUninstall, script+";")
}

func (flufikPkg *FlufikPackage) PreInScript() string {
	return strings.Join(flufikPkg.PreInstall, "\n")
}

func (flufikPkg *FlufikPackage) PostInScript() string {
	return strings.Join(flufikPkg.PostInstall, "\n")
}

func (flufikPkg *FlufikPackage) PreUnScript() string {
	return strings.Join(flufikPkg.PreUninstall, "\n")
}

func (flufikPkg *FlufikPackage) PostUnScript() string {
	return strings.Join(flufikPkg.PreUninstall, "\n")
}

//Signature

func (flufikPkg *FlufikPackage) AddSignatureKey() string {
	return flufikPkg.Signature.PgpName
}

func (flufikPkg *FlufikPackage) AddSignatureType() string {
	return flufikPkg.Signature.Type
}

//func (flufikPkg *FlufikPackage) AddSignaturePassPhrase() string {
//	return flufikPkg.Signature.PassPhrase
//}

func (p *FlufikPackage) JoinedFilePath(filepath string) string {
	//if filepath != "" && !strings.HasPrefix(filepath, "/") {
	//	fmt.Println(path.Join(p.sourceHome, filepath))
	//	return path.Join(p.sourceHome, filepath)
	//} else {
	//	return filepath
	//}
	return filepath
}

func (p *FlufikPackage) init() {
	//p.Logger.Info("Update all package directory which have source directory")
	p.Meta.UpdateBuildTime(time.Now().UTC())

	for k, f := range p.Files {
		for id, file := range f {
			p.Files[k][id].Source = p.JoinedFilePath(file.Source)
		}
	}
}

func (p *FlufikPackage) AddDirectory(pkgdir FlufikPackageDir) {
	p.Directory = append(p.Directory, pkgdir)
}

func (p *FlufikPackage) AddFile(fileT string, pFile FlufikPackageFile) error {
	pFile.Source = p.JoinedFilePath(pFile.Source)

	p.Files[fileT] = append(p.Files[fileT], pFile)

	return nil
}

func LoadPackageInfo(filepath string, sourceHome string) (*FlufikPackage, error) {
	buffer, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	pkg := new(FlufikPackage)

	if err = yaml.Unmarshal(buffer, pkg); err != nil {
		return nil, err
	}

	pkg.sourceHome = sourceHome

	pkg.init()

	return pkg, nil
}
