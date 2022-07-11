package config

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/logging"
	"os"
	"strings"
)

type ServiceConfigBuilder struct {
	RootRepoPath            string   //ENV FLUFIK_ROOT_PATH
	PrivateKey              string   //ENV FLUFIK_PRIVATE_KEY
	PrivateKeyPwd           string   //ENV FLUFIK_PRIVATE_KEY_PWD
	RpmRepositoryName       string   //ENV FLUFIK_RPM_REPO_NAME
	ListenPort              string   //ENV FLUFIK_LPORT
	SupportArch             []string //ENV FLUFIK_SUPPORT_ARCH
	Sections                []string //ENV FLUFIK_SECTIONS
	DistroNames             []string //ENV FLUFIK_DISTRO_NAMES
	EnableSSL               bool     //ENV FLUFIK_ENABLE_SSL
	SSLCert                 string   //base64 ENV FLUFIK_SSL_CERT
	SSLKey                  string   //base64 ENV FLUFIK_SSL_KEY
	EnableSigning           bool     //ENV FLUFIK_ENABLE_SIGNING
	EnableDirectoryWatching bool     //ENV FLUFIK_ENABLE_DIR_WATCH
	PrivateKeyName          string   //ENV FLUFIK_PRIVATE_KEY_NAME
}

func GetServiceConfiguration() *ServiceConfigBuilder {
	var s ServiceConfigBuilder
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
	if os.Getenv("FLUFIK_ENABLE_SSL") == "true" {
		s.EnableSSL = true
		if os.Getenv("FLUFIK_SSL_CERT") == "" || os.Getenv("FLUFIK_SSL_KEY") == "" {
			logging.ErrorHandler("warning: ", fmt.Errorf("Enable SSL enabled but certificate not provided!!!"))
		} else {
			cert, err := crypto.SaveB64DecodedData(os.Getenv("FLUFIK_SSL_CERT"), "server.crt")
			if err != nil {
				logging.ErrorHandler("fatal: ", err)
			}
			key, err := crypto.SaveB64DecodedData(os.Getenv("FLUFIK_SSL_KEY"), "server.key")
			if err != nil {
				logging.ErrorHandler("fatal: ", err)
			}
			s.SSLCert = cert
			s.SSLKey = key
		}
	} else {
		s.EnableSSL = false
	}
	if os.Getenv("FLUFIK_ENABLE_SIGNING") == "true" {
		s.EnableSigning = true
		if os.Getenv("FLUFIK_PRIVATE_KEY_NAME") == "" {
			logging.ErrorHandler("message: ", fmt.Errorf("private key name missing"))
		} else {
			s.PrivateKeyName = os.Getenv("FLUFIK_PRIVATE_KEY_NAME")
		}
	} else {
		s.EnableSigning = false
	}
	if os.Getenv("FLUFIK_RPM_REPO_NAME") == "" {
		s.RpmRepositoryName = "flufik"
	} else {
		s.RpmRepositoryName = os.Getenv("FLUFIK_RPM_REPO_NAME")
	}
	return &s
}
