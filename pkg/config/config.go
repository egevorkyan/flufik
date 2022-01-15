package config

import (
	"fmt"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/logging"
	"os"
	"strings"
)

type ServiceConfigBuilder struct {
	ListenPort              string   //ENV FLUFIK_LPORT
	RootRepoPath            string   //ENV FLUFIK_ROOT_PATH
	SupportArch             []string //ENV FLUFIK_SUPPORT_ARCH
	Sections                []string //ENV FLUFIK_SECTIONS
	DistroNames             []string //ENV FLUFIK_DISTRO_NAMES
	EnableSSL               bool     //ENV FLUFIK_ENABLE_SSL
	SSLCert                 string   //base64 ENV FLUFIK_SSL_CERT
	SSLKey                  string   //base64 ENV FLUFIK_SSL_KEY
	EnableAPIKeys           bool     //ENV FLUFIK_ENABLE_API_KEYS
	EnableSigning           bool     //ENV FLUFIK_ENABLE_SIGNING
	PrivateKey              string   //base64 ENV FLUFIK_PRIVATE_KEY
	PrivateKeyPwd           string   //base64 ENV FLUFIK_PRIVATE_KEY_PWD
	EnableDirectoryWatching bool     //ENV FLUFIK_ENABLE_DIR_WATCH
}

func GetServiceConfiguration() *ServiceConfigBuilder {
	var s ServiceConfigBuilder
	if os.Getenv("FLUFIK_LPORT") == "" {
		s.ListenPort = "8080"
	} else {
		s.ListenPort = os.Getenv("FLUFIK_LPORT")
	}
	if os.Getenv("FLUFIK_ROOT_PATH") == "" {
		s.RootRepoPath = "opt/flufik"
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
	if os.Getenv("FLUFIK_ENABLE_API_KEYS") == "true" {
		s.EnableAPIKeys = true
	} else {
		s.EnableAPIKeys = false
	}
	if os.Getenv("FLUFIK_ENABLE_SIGNING") == "true" {
		s.EnableSigning = true
		if os.Getenv("FLUFIK_PRIVATE_KEY") == "" || os.Getenv("FLUFIK_PRIVATE_KEY_PWD") == "" {
			//_ := crypto.PasswordGenerator(10, 2, 3, 4)

		}
	}
	return &s
}
