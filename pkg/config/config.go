package config

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"os"
	"strings"
)

type ServiceConfigBuilder struct {
	PublicUrl                  string   //ENV FLUFIK_PUBLIC_URL
	RootRepoPath               string   //ENV FLUFIK_ROOT_PATH
	RpmRepositoryName          string   //ENV FLUFIK_RPM_REPO_NAME
	RpmRepositoryOsName        []string //ENV FLUFIK_RPM_REPO_SUPPORTED_OSNAME
	RpmRepositoryRhelVersion   []string //ENV FLUFIK_RPM_REPO_RHEL_SUPPORTED_VERSION
	RpmRepositoryFedoraVersion []string //ENV FLUFIK_RPM_REPO_FEDORA_SUPPORTED_VERSION
	RpmRepositoryArch          []string //ENV FLUFIK_RPM_REPO_SUPPORTED_ARCH
	ListenPort                 string   //ENV FLUFIK_LPORT
	SupportArch                []string //ENV FLUFIK_SUPPORT_ARCH
	Sections                   []string //ENV FLUFIK_SECTIONS
	DistroNames                []string //ENV FLUFIK_DISTRO_NAMES
}

func GetServiceConfiguration() (*ServiceConfigBuilder, error) {
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
	if os.Getenv("FLUFIK_RPM_REPO_SUPPORTED_OSNAME") == "" {
		s.RpmRepositoryOsName = []string{"rhel", "centos", "fedora"}
	} else {
		s.RpmRepositoryOsName = strings.Split(os.Getenv("FLUFIK_RPM_REPO_SUPPORTED_OSNAME"), " ")
	}
	if os.Getenv("FLUFIK_RPM_REPO_RHEL_SUPPORTED_VERSION") == "" {
		s.RpmRepositoryRhelVersion = []string{"7", "8", "9"}
	} else {
		s.RpmRepositoryRhelVersion = strings.Split(os.Getenv("FLUFIK_RPM_REPO_RHEL_SUPPORTED_VERSION"), " ")
	}
	if os.Getenv("FLUFIK_RPM_REPO_FEDORA_SUPPORTED_VERSION") == "" {
		s.RpmRepositoryFedoraVersion = []string{"7", "8", "9"}
	} else {
		s.RpmRepositoryFedoraVersion = strings.Split(os.Getenv("FLUFIK_RPM_REPO_FEDORA_SUPPORTED_VERSION"), " ")
	}
	if os.Getenv("FLUFIK_RPM_REPO_SUPPORTED_ARCH") == "" {
		s.RpmRepositoryArch = []string{"noarch", "aarch64", "x86_64", "s390x"}
	} else {
		s.RpmRepositoryArch = strings.Split(os.Getenv("FLUFIK_RPM_REPO_SUPPORTED_ARCH"), " ")
	}
	return &s, nil
}
