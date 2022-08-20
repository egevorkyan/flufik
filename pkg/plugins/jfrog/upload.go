package jfrog

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/logging"
	"io"
	"net/http"
)

type JFrog struct {
	repoUser     string
	repoPwd      string
	repoUrl      string
	packageName  string
	path         string
	distribution string
	component    string
	architecture string
	repoName     string
	logger       *logging.Logger
	debugging    string
}

func (j *JFrog) FlufikJFrogUpload() error {
	if j.debugging == "1" {
		j.logger.Info("uploading to jfrog repository")
	}
	var requestUrl string
	pkgFile, err := core.OpenFile(j.packageName, j.path)
	if err != nil {
		return fmt.Errorf("can not open package: %v", err)
	}
	pkg := core.CheckPackage(j.packageName)
	if pkg == "deb" {
		requestUrl = fmt.Sprintf("%s/artifactory/%s/pool/%s;deb.distribution=%s;deb.component=%s;deb.architecture=%s", j.repoUrl, j.repoName,
			j.packageName, j.distribution, j.component, j.architecture)
	} else if pkg == "rpm" {
		requestUrl = fmt.Sprintf("%s/artifactory/%s/%s/%s", j.repoUrl, j.repoName, j.distribution, j.packageName)
	} else {
		return fmt.Errorf("failure: %s", pkg)
	}
	h := crypto.NewHash(pkgFile.Name(), j.logger, j.debugging)
	checksum, err := h.CheckSum()
	if err != nil {
		return fmt.Errorf("no checksum %v", err)
	}
	req, err := http.NewRequest("PUT", requestUrl, pkgFile)
	if err != nil {
		return fmt.Errorf("request builder failure: %w", err)
	}
	req.SetBasicAuth(j.repoUser, j.repoPwd)

	req.Header.Set("X-Checksum-Sha1", checksum.Sha1)
	req.Header.Set("X-Checksum-Sha256", checksum.Sha256)
	req.Header.Set("X-Checksum-Md5", checksum.Md5)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("response failed: %v", err)
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("upload failed: %v", response.Status)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	return nil
}

func NewUpload(repoUser, repoPwd, repoUrl, packageName, path, distribution, component, architecture, repoName string, logger *logging.Logger, debugging string) *JFrog {
	j := &JFrog{
		repoUser:     repoUser,
		repoPwd:      repoPwd,
		repoUrl:      repoUrl,
		packageName:  packageName,
		path:         path,
		distribution: distribution,
		component:    component,
		architecture: architecture,
		repoName:     repoName,
		logger:       logger,
		debugging:    debugging,
	}
	return j
}
