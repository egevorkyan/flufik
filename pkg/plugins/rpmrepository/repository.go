package rpmrepository

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/config"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/jfrog/go-rpm"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const repoDataPath = "repodata"

type RpmRepo struct {
	cfg       *config.ServiceConfigBuilder
	logger    *logging.Logger
	debugging string
}

var packageinfo = map[string]PackageInfos{}
var repodata = map[string]RepoData{}

func NewRpmBuilder(config *config.ServiceConfigBuilder, logger *logging.Logger, debugging string) *RpmRepo {
	return &RpmRepo{cfg: config, logger: logger, debugging: debugging}
}

// CreateBaseDir - creates base directory
func (r *RpmRepo) CreateBaseDir() error {
	if r.debugging == "1" {
		r.logger.Info("create rpm repository base directory")
	}
	fmt.Println(filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName))
	err := r.createDir(filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName))
	if err != nil {
		return fmt.Errorf("failed to create base directory: %v", err)
	}
	return nil
}

// ReindexPackages - reindex packages on repository server
func (r *RpmRepo) ReindexPackages() error {
	if r.debugging == "1" {
		r.logger.Info("reindexing packages")
	}
	// reset settings
	packageinfo = map[string]PackageInfos{}
	repodata = map[string]RepoData{}

	// find repos
	elems, err := ioutil.ReadDir(filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName))
	if err != nil {
		return fmt.Errorf("failure occur during reading packages: %v", err)
	}
	for _, elem := range elems {
		if elem.IsDir() {
			packageinfo[elem.Name()] = PackageInfos{}
		}
	}
	// read rpms in repo
	for repo, _ := range packageinfo {
		err = filepath.Walk(filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName, repo), func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				if strings.HasSuffix(f.Name(), "rpm") {
					// get sha256
					file, err := os.Open(path)
					if err != nil {
						return fmt.Errorf("can not open file: %v", err)
					}
					defer func(file *os.File) {
						err = file.Close()
						if err != nil {
							return
						}
					}(file)
					hasher := sha256.New()
					if _, err := io.Copy(hasher, file); err != nil {
						return fmt.Errorf("can not copy file to hasher: %v", err)
					}
					sumString := fmt.Sprintf("%x", hasher.Sum(nil))
					// get rpm info
					p, err := rpm.OpenPackageFile(path)
					if err != nil {
						return fmt.Errorf("can not open rpm file: %v", err)
					}
					pi := PackageInfo{f.Name(), *p}
					// store
					packageinfo[repo][sumString] = pi
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		repodata[repo], err = r.CreateRepoData(packageinfo[repo])
		if err != nil {
			return fmt.Errorf("can not create repodata: %v", err)
		}
	}
	return nil
}

// SavePackage - saves uploaded package in repository
func (r *RpmRepo) SavePackage(data io.Reader, pkg string) error {
	if r.debugging == "1" {
		r.logger.Info("saving package")
	}
	filePath := filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName, pkg)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("can not create file %v: %v", filePath, err)
	}
	defer func(out *os.File) {
		err = out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, data)
	if err != nil {
		return fmt.Errorf("can not copy package: %v", err)
	}
	return nil
}

// readRepos - reading repositories
func (r *RpmRepo) readRepos() error {
	return nil
}

// getPaths - building paths and returns paths
func (r *RpmRepo) getPaths(supportedOs string) []string {
	if r.debugging == "1" {
		r.logger.Info("building path structure and returning value")
	}
	var paths []string
	supported := strings.Split(supportedOs, " ")
	for _, s := range supported {
		temp := filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName, s)
		paths = append(paths, temp)
	}
	return paths
}

// createDir - creates directories
func (r *RpmRepo) createDir(path string) error {
	if r.debugging == "1" {
		r.logger.Info("create directory")
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	return nil
}

func (r *RpmRepo) Repository(uploaded []byte, uploadedName string) error {
	if r.debugging == "1" {
		r.logger.Info("rpm repository")
	}
	sum := sha256.Sum256(uploaded)
	sumString := fmt.Sprintf("%x", sum)
	p, err := rpm.ReadPackageFile(bytes.NewBuffer(uploaded))
	if err != nil {
		return fmt.Errorf("failed reading rpm file: %v", err)
	}
	pi := PackageInfo{p.String(), *p}
	err = r.SavePackage(bytes.NewBuffer(uploaded), uploadedName)
	if err != nil {
		return err
	}
	temp := make(PackageInfos)
	temp[sumString] = pi
	packageinfo[r.cfg.RpmRepositoryName] = temp
	repodata[r.cfg.RpmRepositoryName], err = r.CreateRepoData(packageinfo[r.cfg.RpmRepositoryName])
	if err != nil {
		return err
	}
	err = r.createDir(filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName, repoDataPath))
	if err != nil {
		return err
	}
	err = r.Dump(repodata[r.cfg.RpmRepositoryName])
	if err != nil {
		return err
	}
	return nil
}

func (r *RpmRepo) Dump(data RepoData) error {
	if r.debugging == "1" {
		r.logger.Info("dumping repository xml files")
	}
	for k, v := range data {
		if strings.Contains(k, "gz") || strings.Contains(k, "repomd.xml") {
			f, err := os.Create(filepath.Join(core.FlufikServiceWebHome(r.cfg.RootRepoPath), r.cfg.RpmRepositoryName, repoDataPath, k))
			if err != nil {
				return fmt.Errorf("can not create file: %v", err)
			}
			defer func(f *os.File) {
				err = f.Close()
				if err != nil {
					return
				}
			}(f)
			_, err = f.Write(v)
			if err != nil {
				return fmt.Errorf("can not write to file: %v", err)
			}
		}
	}
	return nil
}
