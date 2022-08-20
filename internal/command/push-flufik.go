package command

import (
	"bytes"
	"fmt"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type PushRepoFlufikCommand struct {
	command  *cobra.Command
	pkgType  string // -t
	arch     string // -a
	distro   string // -d
	section  string // -s
	repoUrl  string // -f
	username string // -u
	password string // -p
	pkg      string // -b
}

func NewFlufikPushRepoFlufikCommand() *PushRepoFlufikCommand {
	c := &PushRepoFlufikCommand{
		command: &cobra.Command{
			Use:   "flufik",
			Short: "push package to flufik repository",
		},
	}
	c.command.Flags().StringVarP(&c.repoUrl, "url", "f", "", "flufik.com")
	c.command.Flags().StringVarP(&c.username, "username", "u", "", "authentication username")
	c.command.Flags().StringVarP(&c.password, "password", "p", "", "authentication password")
	c.command.Flags().StringVarP(&c.pkgType, "package-type", "t", "rpm", "deb or rpm packages")
	c.command.Flags().StringVarP(&c.arch, "arch", "a", "amd64", "Required: for debian repository only. CPU architecture")
	c.command.Flags().StringVarP(&c.section, "section", "s", "main", "Required: for debian repository only")
	c.command.Flags().StringVarP(&c.distro, "distro", "d", "", "Required: for debian repository only")
	c.command.Flags().StringVarP(&c.pkg, "pkg-file", "b", "", "package path")
	c.command.Run = c.Run
	_ = c.command.MarkFlagRequired("url")
	_ = c.command.MarkFlagRequired("package-type")
	_ = c.command.MarkFlagRequired("pkg-file")
	_ = c.command.MarkFlagRequired("username")
	_ = c.command.MarkFlagRequired("password")
	return c
}

func (c *PushRepoFlufikCommand) Run(command *cobra.Command, args []string) {
	logger := logging.GetLogger()
	debuging := os.Getenv("FLUFIK_DEBUG")
	if debuging == "1" {
		logger.Info("publishing package")
	}
	extraParams := make(map[string]string)
	switch c.pkgType {
	case "deb":
		if c.distro == "" && c.pkgType == "" && c.arch == "" {
			logger.Error("missing parameter values")
		}
		urlBuilder := fmt.Sprintf("https://%s/upload?type=%s&arch=%s&distro=%s&section=%s", c.repoUrl, c.pkgType, c.arch, c.distro, c.section)
		request, err := newfileUploadRequest(urlBuilder, extraParams, "file", c.pkg, c.username, c.password)
		if err != nil {
			logger.Errorf(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			logger.Errorf("failed during upload request: %v", err)
		} else {
			body := &bytes.Buffer{}
			_, err = body.ReadFrom(resp.Body)
			if err != nil {
				logger.Errorf("failed to read response body: %v", err)
			}
			err = resp.Body.Close()
			if err != nil {
				logger.Errorf("can not close body: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				logger.Fatalf("failed to upload file: %s", c.pkg)
			}
			logger.Infof("Status: %s, StatusCode: %v", resp.Status, resp.StatusCode)
		}

	case "rpm":
		urlBuilder := fmt.Sprintf("https://%s/upload?type=%s&arch=%s&distro=%s&section=%s", c.repoUrl, c.pkgType, c.arch, c.distro, c.section)
		request, err := newfileUploadRequest(urlBuilder, extraParams, "file", c.pkg, c.username, c.password)
		if err != nil {
			logger.Errorf(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			logger.Errorf("failed during upload request: %v", err)
		} else {
			body := &bytes.Buffer{}
			_, err = body.ReadFrom(resp.Body)
			if err != nil {
				logger.Errorf("failed to read response body: %v", err)
			}
			err = resp.Body.Close()
			if err != nil {
				logger.Errorf("can not close body: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				logger.Fatalf("failed to upload file: %s", c.pkg)
			}
			logger.Infof("Status: %s, StatusCode: %v", resp.Status, resp.StatusCode)
		}
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string, username, password string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
