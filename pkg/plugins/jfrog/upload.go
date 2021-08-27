package jfrog

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
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
}

func (j *JFrog) FlufikJFrogUpload() error {
	var requestUrl string
	pkgFile, err := core.OpenFile(j.packageName, j.path)
	if err != nil {
		logging.ErrorHandler("can not open package: ", err)
	}
	pkg := core.CheckPackage(j.packageName)
	if pkg == "deb" {
		requestUrl = fmt.Sprintf("%s/%s;deb.distribution=%s;deb.component=%s;deb.architecture=%s", j.repoUrl,
			j.packageName, j.distribution, j.component, j.architecture)
	} else if pkg == "rpm" {
		requestUrl = fmt.Sprintf("%s/%s/%s", j.repoUrl, j.distribution, j.packageName)
	} else {
		return fmt.Errorf("failure: %s", pkg)
	}
	checksum, err := core.CheckSum(pkgFile.Name())
	if err != nil {
		return fmt.Errorf("no checksum", err)
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
		return fmt.Errorf("response failed: %w", err)
	}

	fmt.Println(response.Status, response.StatusCode)

	if response.StatusCode != 201 {
		return fmt.Errorf("upload failed: %s", response.Status)
	}

	defer response.Body.Close()

	return nil
}

func NewUpload(repoUser, repoPwd, repoUrl, packageName, path, distribution, component, architecture string) *JFrog {
	j := &JFrog{
		repoUser:     repoUser,
		repoPwd:      repoPwd,
		repoUrl:      repoUrl,
		packageName:  packageName,
		path:         path,
		distribution: distribution,
		component:    component,
		architecture: architecture,
	}
	return j
}
