package config

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"os"
	"strings"
)

type ServiceConfigBuilder struct {
	PublicUrl         string   //ENV FLUFIK_PUBLIC_URL
	RootRepoPath      string   //ENV FLUFIK_ROOT_PATH
	RpmRepositoryName string   //ENV FLUFIK_RPM_REPO_NAME
	ListenPort        string   //ENV FLUFIK_LPORT
	SupportArch       []string //ENV FLUFIK_SUPPORT_ARCH
	Sections          []string //ENV FLUFIK_SECTIONS
	DistroNames       []string //ENV FLUFIK_DISTRO_NAMES
}

func GetServiceConfiguration(logger *logging.Logger, debugging string) (*ServiceConfigBuilder, error) {
	if debugging == "1" {
		logger.Info("service configuration file")
	}
	var s ServiceConfigBuilder
	if os.Getenv("FLUFIK_PUBLIC_URL") == "" {
		return nil, fmt.Errorf("public url is missing")
	} else {
		s.PublicUrl = os.Getenv("FLUFIK_PUBLIC_URL")
	}
	if os.Getenv("FLUFIK_LPORT") == "" {
		s.ListenPort = "8080"
	} else {
		s.ListenPort = os.Getenv("FLUFIK_LPORT")
	}
	if os.Getenv("FLUFIK_ROOT_PATH") == "" {
		s.RootRepoPath = core.FlufikRootHome()
	} else {
		s.RootRepoPath = os.Getenv("FLUFIK_ROOT_PATH")
	}
	if os.Getenv("FLUFIK_SUPPORT_ARCH") == "" {
		s.SupportArch = []string{"386", "amd64", "arm64"}
	} else {
		s.SupportArch = strings.Split(os.Getenv("FLUFIK_SUPPORT_ARCH"), " ")
	}
	if os.Getenv("FLUFIK_SECTIONS") == "" {
		s.Sections = []string{}
	} else {
		s.Sections = strings.Split(os.Getenv("FLUFIK_SECTIONS"), " ")
	}
	if os.Getenv("FLUFIK_DISTRO_NAMES") == "" {
		s.DistroNames = []string{}
	} else {
		s.DistroNames = strings.Split(os.Getenv("FLUFIK_DISTRO_NAMES"), " ")
	}
	if os.Getenv("FLUFIK_RPM_REPO_NAME") == "" {
		s.RpmRepositoryName = "flufik"
	} else {
		s.RpmRepositoryName = os.Getenv("FLUFIK_RPM_REPO_NAME")
	}
	return &s, nil
}
