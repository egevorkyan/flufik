package command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type PushRepoFlufikCommand struct {
	command   *cobra.Command
	arch      string // -a
	distro    string // -d
	section   string // -s
	repoUrl   string // -f
	username  string // -u
	password  string // -p
	pkg       string // -b
	osversion string // -o
	isYum     bool
	isApt     bool
}

func NewFlufikPushRepoFlufikCommand() *PushRepoFlufikCommand {
	c := &PushRepoFlufikCommand{
		command: &cobra.Command{
			Use:   "generic",
			Short: "push package to flufik repository",
		},
	}
	c.command.Flags().StringVarP(&c.repoUrl, "url", "f", "", "flufik.com")
	c.command.Flags().StringVarP(&c.username, "username", "u", "", "authentication username")
	c.command.Flags().StringVarP(&c.password, "password", "p", "", "authentication password")
	c.command.Flags().StringVarP(&c.arch, "arch", "a", "amd64", "cpu architecture")
	c.command.Flags().StringVarP(&c.section, "section", "s", "main", "Required: for debian repository only")
	c.command.Flags().StringVarP(&c.distro, "distro", "d", "", "os distribution")
	c.command.Flags().StringVarP(&c.pkg, "pkg-file", "b", "", "package path")
	c.command.Flags().BoolVarP(&c.isYum, "yum", "", false, "indicates pushing yum based package")
	c.command.Flags().BoolVarP(&c.isApt, "apt", "", false, "indicates pushing apt based package")
	c.command.Flags().StringVarP(&c.osversion, "os-version", "", "", "os version for rhel based os")
	c.command.Run = c.Run
	_ = c.command.MarkFlagRequired("url")
	_ = c.command.MarkFlagRequired("distro")
	_ = c.command.MarkFlagRequired("pkg-file")
	_ = c.command.MarkFlagRequired("username")
	_ = c.command.MarkFlagRequired("password")
	return c
}

func (c *PushRepoFlufikCommand) Run(command *cobra.Command, args []string) {
	var arch string
	extraParams := make(map[string]string)
	if c.isApt {
		if c.distro == "" && c.arch == "" {
			logger.RaiseErr("missing parameter values")
		}
		if c.arch == "x86_64" {
			arch = "amd64"
		} else if c.arch == "aarch64" {
			arch = "arm64"
		} else {
			arch = c.arch
		}
		urlBuilder := fmt.Sprintf("https://%s/upload/apt?&arch=%s&distro=%s&section=%s", c.repoUrl, arch, c.distro, c.section)
		request, err := newfileUploadRequest(urlBuilder, extraParams, "file", c.pkg, c.username, c.password)
		if err != nil {
			logger.RaiseErr("failed request", err)
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			logger.RaiseErr("failed during upload request", err)
		} else {
			body := &bytes.Buffer{}
			_, err = body.ReadFrom(resp.Body)
			if err != nil {
				logger.RaiseErr("failed to read response body", err)
			}
			err = resp.Body.Close()
			if err != nil {
				logger.RaiseErr("can not close body", err)
			}
			if resp.StatusCode != http.StatusOK {
				logger.RaiseErr("failed to upload file", errors.New(c.pkg))
			}
			logger.InfoLog("Status: %s, StatusCode: %v", resp.Status, resp.StatusCode)
		}
	} else if c.isYum {
		if c.distro == "" && c.arch == "" && c.osversion == "" {
			logger.RaiseErr("missing parameter values")
		}
		if c.arch == "amd64" {
			arch = "x86_64"
		} else if c.arch == "arm64" {
			arch = "aarch64"
		} else {
			arch = c.arch
		}
		urlBuilder := fmt.Sprintf("https://%s/upload/yum?&arch=%s&distro=%s&version=%s", c.repoUrl, arch, c.distro, c.osversion)
		request, err := newfileUploadRequest(urlBuilder, extraParams, "file", c.pkg, c.username, c.password)
		if err != nil {
			logger.RaiseErr("failed request", err)
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			logger.RaiseErr("failed during upload request", err)
		} else {
			body := &bytes.Buffer{}
			_, err = body.ReadFrom(resp.Body)
			if err != nil {
				logger.RaiseErr("failed to read response body", err)
			}
			err = resp.Body.Close()
			if err != nil {
				logger.RaiseErr("can not close body", err)
			}
			if resp.StatusCode != http.StatusOK {
				logger.RaiseErr("failed to upload file", errors.New(c.pkg))
			}
			logger.InfoLog("Status: %s, StatusCode: %v", resp.Status, resp.StatusCode)
		}
	} else {
		logger.RaiseWarn("Required flags are missing")
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
