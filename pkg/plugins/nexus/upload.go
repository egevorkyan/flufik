package nexus

import (
	"bytes"
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type FlufikNexus struct {
	repoUser     string
	repoPwd      string
	repoUrl      string
	pkgName      string
	path         string
	nxcomponent  string
	nxrepository string
}

func (fn *FlufikNexus) FlufikNexusUpload() error {
	requestUrl := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", fn.repoUrl, fn.nxrepository)
	pkg, err := os.Open(core.FlufikPkgFilePath(fn.pkgName, fn.path))
	if err != nil {
		return err
	}
	pkgType := core.CheckPackage(fn.pkgName)

	if pkgType == "deb" {
		if err = fn.debUpload(pkg, requestUrl); err != nil {
			return err
		}
	} else if pkgType == "rpm" {
		if err = fn.rpmUpload(pkg, requestUrl); err != nil {
			return err
		}
	}
	return nil
}

func (fn *FlufikNexus) debUpload(pkg *os.File, requestUrl string) error {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	mpart, err := w.CreateFormFile("apt.asset", fn.pkgName)
	if err != nil {
		return err
	}
	_, err = io.Copy(mpart, pkg)
	err = w.Close()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", requestUrl, body)
	if err != nil {
		return fmt.Errorf("request builder failure: %w", err)
	}

	req.SetBasicAuth(fn.repoUser, fn.repoPwd)
	req.Header.Set("Content-Type", w.FormDataContentType())

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("response failed: %w", err)
	}

	if response.StatusCode != 204 {
		return fmt.Errorf("upload failed: %s", response.Status)
	}

	defer response.Body.Close()

	return nil
}

func (fn *FlufikNexus) rpmUpload(pkg *os.File, requestUrl string) error {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	mpart, err := w.CreateFormFile("yum.asset", fn.pkgName)
	if err != nil {
		return err
	}
	_, err = io.Copy(mpart, pkg)

	if err = w.WriteField("yum.asset.filename", fn.pkgName); err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", requestUrl, body)
	if err != nil {
		return fmt.Errorf("request builder failure: %w", err)
	}

	req.SetBasicAuth(fn.repoUser, fn.repoPwd)
	req.Header.Set("Content-Type", w.FormDataContentType())

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("response failed: %w", err)
	}

	if response.StatusCode != 204 {
		return fmt.Errorf("upload failed: %s", response.Status)
	}

	defer response.Body.Close()

	return nil
}

func NewNexusUpload(repoUser, repoPwd, repoUrl, packageName, path, nxcomponent, nxrepository string) *FlufikNexus {
	n := &FlufikNexus{
		repoUser:     repoUser,
		repoPwd:      repoPwd,
		repoUrl:      repoUrl,
		pkgName:      packageName,
		path:         path,
		nxcomponent:  nxcomponent,
		nxrepository: nxrepository,
	}
	return n
}
