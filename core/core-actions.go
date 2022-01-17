package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func OpenFile(fileName string, path string) (*os.File, error) {
	if !filepath.IsAbs(path) {
		p, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("can't make path absolute: %w", err)
		}
		path = p
	}
	fullPath := filepath.Join(path, fileName)
	pkg, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("can not open file: %w", err)
	}
	return pkg, nil
}

func CheckPackage(fileName string) string {
	if strings.HasSuffix(fileName, ".deb") {
		return "deb"
	} else if strings.HasSuffix(fileName, ".rpm") {
		return "rpm"
	} else {
		return "unknown file extension"
	}
}

func CheckArch(packageName string) string {
	ext := CheckPackage(packageName)
	var arch string
	if ext == "deb" {
		if strings.Contains(packageName, "i386") {
			arch = "i386"
		} else if strings.Contains(packageName, "amd64") {
			arch = "amd64"
		} else if strings.Contains(packageName, "armhf") {
			arch = "armhf"
		} else if strings.Contains(packageName, "arm64") {
			arch = "arm64"
		} else if strings.Contains(packageName, "ppc64el") {
			arch = "ppc64el"
		} else if strings.Contains(packageName, "s390x") {
			arch = "s390x"
		}
	} else if ext == "rpm" {
		if strings.Contains(packageName, "x86_64") {
			arch = "x86_64"
		} else if strings.Contains(packageName, "aarch64") {
			arch = "aarch64"
		} else if strings.Contains(packageName, "s390x") {
			arch = "s390x"
		} else if strings.Contains(packageName, "s390") {
			arch = "s390"
		} else if strings.Contains(packageName, "ia64") {
			arch = "ia64"
		} else if strings.Contains(packageName, "s390x") {
			arch = "s390x"
		} else if strings.Contains(packageName, "ppc64le") {
			arch = "ppc64le"
		} else if strings.Contains(packageName, "noarch") {
			arch = "noarch"
		}
	}
	return arch
}
