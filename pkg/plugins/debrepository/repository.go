package debrepository

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/blakesmith/ar"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/fsnotify/fsnotify"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
)

type Compression int

const (
	LZMA Compression = iota
	GZIP
)

func (s *ServiceConfigBuilder) ArchPath(distro string, section string, arch string) string {
	return filepath.Join(core.FlufikServiceWebHome(), "dists", distro, section, "binary-"+arch)
}

func (s *ServiceConfigBuilder) CreateDirectories() error {
	for _, distro := range s.DistroNames {
		for _, arch := range s.SupportArch {
			for _, section := range s.Sections {
				if _, err := os.Stat(s.ArchPath(distro, section, arch)); err != nil {
					if os.IsNotExist(err) {
						if err = os.MkdirAll(s.ArchPath(distro, section, arch), 0755); err != nil {
							return err
						}
					} else {
						return fmt.Errorf("error inspecting %s (%s): %s", distro, arch, err)
					}
				}
				if s.EnableDirectoryWatching {
					if err := s.directoryWatcher.Add(s.ArchPath(distro, section, arch)); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (s *ServiceConfigBuilder) RebuildRepoMetadata(filePath string) error {
	distroArch := destructPath(filePath)
	if err := s.createPackageGz(distroArch[0], distroArch[1], distroArch[2]); err != nil {
		return err
	}
	if s.EnableSigning {
		if err := s.createRelease(distroArch[0]); err != nil {
			return err
		}

	}
	return nil
}

func (s *ServiceConfigBuilder) createRelease(distro string) error {
	workingDirectory := filepath.Join(core.FlufikServiceWebHome(), "dists", distro)

	outFile, err := os.Create(filepath.Join(workingDirectory, "Release"))
	if err != nil {
		return fmt.Errorf("failed to create Release: %s", err)
	}
	defer outFile.Close()

	currentTime := time.Now().UTC()
	fmt.Fprintf(outFile, "Suite: %s\n", distro)
	fmt.Fprintf(outFile, "Codename: %s\n", distro)
	fmt.Fprintf(outFile, "Components: %s\n", strings.Join(s.Sections, " "))
	fmt.Fprintf(outFile, "Architectures: %s\n", strings.Join(s.SupportArch, " "))
	fmt.Fprintf(outFile, "Date: %s\n", currentTime.Format("Mon, 02 Jan 2006 15:04:05 UTC"))

	var md5Sums strings.Builder
	var sha1Sums strings.Builder
	var sha256Sums strings.Builder

	err = filepath.Walk(workingDirectory, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "Package.gz") || strings.HasSuffix(path, "Packages") {
			var (
				md5hash    = md5.New()
				sha1hash   = sha1.New()
				sha256hash = sha256.New()
			)
			relPath, _ := filepath.Rel(workingDirectory, path)
			slashPath := filepath.ToSlash(relPath)
			f, err := os.Open(path)
			if err != nil {
				log.Println("error opening the packages file for reading", err)
			}
			if _, err = io.Copy(io.MultiWriter(md5hash, sha1hash, sha256hash), f); err != nil {
				return fmt.Errorf("error hashing file for release list: %s", err)
			}
			fmt.Fprintf(&md5Sums, " %s %d %s\n", hex.EncodeToString(md5hash.Sum(nil)), file.Size(), slashPath)
			fmt.Fprintf(&sha1Sums, " %s %d %s\n", hex.EncodeToString(sha1hash.Sum(nil)), file.Size(), slashPath)
			fmt.Fprintf(&sha256Sums, " %s %d %s\n", hex.EncodeToString(sha256hash.Sum(nil)), file.Size(), slashPath)

			f = nil
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error scaning for packages files: %s", err)
	}

	outFile.WriteString("MD5Sum:\n")
	outFile.WriteString(md5Sums.String())
	outFile.WriteString("SHA1:\n")
	outFile.WriteString(sha1Sums.String())
	outFile.WriteString("SHA256:\n")
	outFile.WriteString(sha256Sums.String())

	if err = crypto.SignRelease(s.PrivateKeyName, outFile.Name()); err != nil {
		return fmt.Errorf("error signing release file: %s", err)
	}

	return nil
}

func (s *ServiceConfigBuilder) createPackageGz(distro string, section string, arch string) error {
	packageFile, err := os.Create(filepath.Join(s.ArchPath(distro, section, arch), "Packages"))
	if err != nil {
		return fmt.Errorf("failed to create Packages: %s", err)
	}
	packageGzFile, err := os.Create(filepath.Join(s.ArchPath(distro, section, arch), "Packages.gz"))
	if err != nil {
		return fmt.Errorf("failed to create Packages.gz: %s", err)
	}
	defer packageFile.Close()
	defer packageGzFile.Close()
	gzOut := gzip.NewWriter(packageGzFile)
	defer gzOut.Close()

	writer := io.MultiWriter(packageFile, gzOut)

	dirList, err := ioutil.ReadDir(s.ArchPath(distro, section, arch))
	if err != nil {
		return fmt.Errorf("scanning: %s: %s", s.ArchPath(distro, section, arch), err)
	}
	for i, debFile := range dirList {
		if strings.HasSuffix(debFile.Name(), "deb") {
			var packageBuffer bytes.Buffer
			debPath := filepath.Join(s.ArchPath(distro, section, arch), debFile.Name())
			tempControlData, err := inspectPackage(debPath)
			if err != nil {
				return err
			}
			packageBuffer.WriteString(tempControlData)
			dir := filepath.ToSlash(filepath.Join("dists", distro, section, "binary-"+arch, debFile.Name()))
			fmt.Fprintf(&packageBuffer, "Filename: %s\n", dir)
			fmt.Fprintf(&packageBuffer, "Size: %d\n", debFile.Size())
			f, err := os.Open(debPath)
			if err != nil {
				log.Println("error opening deb file: ", err)
			}
			defer f.Close()

			var (
				md5hash    = md5.New()
				sha1hash   = sha1.New()
				sha256hash = sha256.New()
			)
			_, err = io.Copy(io.MultiWriter(md5hash, sha1hash, sha256hash), f)
			if err != nil {
				return fmt.Errorf("error hashing file for packages file: %s", err)
			}
			fmt.Fprintf(&packageBuffer, "MD5sum: %s\n", hex.EncodeToString(md5hash.Sum(nil)))
			fmt.Fprintf(&packageBuffer, "SHA1: %s\n", hex.EncodeToString(sha1hash.Sum(nil)))
			fmt.Fprintf(&packageBuffer, "SHA256: %s\n", hex.EncodeToString(sha256hash.Sum(nil)))
			if i != (len(dirList) - 1) {
				packageBuffer.WriteString("\n\n")
			}
			writer.Write(packageBuffer.Bytes())
			f = nil
		}
	}
	return nil
}

func destructPath(filePath string) []string {
	splitPath := strings.Split(filePath, "/")
	archFull := splitPath[len(splitPath)-2]
	archSplit := strings.Split(archFull, "-")
	distro := splitPath[len(splitPath)-4]
	section := splitPath[len(splitPath)-3]
	return []string{distro, section, archSplit[1]}
}

func inspectPackageControl(compression Compression, fileName bytes.Buffer) (string, error) {
	var tarReader *tar.Reader
	var err error
	switch compression {
	case GZIP:
		var compressedFile *gzip.Reader
		compressedFile, err = gzip.NewReader(bytes.NewReader(fileName.Bytes()))
		tarReader = tar.NewReader(compressedFile)
		break
	case LZMA:
		var compressedFile *lzma.Reader
		compressedFile, err = lzma.NewReader(bytes.NewReader(fileName.Bytes()))
		tarReader = tar.NewReader(compressedFile)
		break
	}

	var controlBuffer bytes.Buffer
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to inspect package: %s", err)
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			switch name {
			case "control", "./control":
				io.Copy(&controlBuffer, tarReader)
				return strings.TrimRight(controlBuffer.String(), "\n") + "\n", nil
			}
		default:
			log.Printf("Unable to figure out type : %c in file %s\n", header.Typeflag, name)
		}
	}
	err = nil
	return "", err
}

func inspectPackage(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("error opening package file %s: %s", fileName, err)
	}

	arReader := ar.NewReader(f)
	defer f.Close()
	var controlBuffer bytes.Buffer

	for {
		header, err := arReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error in inspectPackage loop: %s", err)
		}
		if strings.Contains(header.Name, "control.tar") {
			var compression Compression
			if strings.TrimRight(header.Name, "/") == "control.tar.gz" {
				compression = GZIP
			} else if strings.TrimRight(header.Name, "/") == "control.tar.xz" {
				compression = LZMA
			} else {
				err := errors.New("No control file found")
				return "", err
			}
			io.Copy(&controlBuffer, arReader)
			return inspectPackageControl(compression, controlBuffer)
		}
	}
	return "", nil
}

func (s *ServiceConfigBuilder) DirectoryWatch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("error creating fswatcher: ", err)
		return err
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Remove == fsnotify.Remove) {
					mutex.Lock()
					if filepath.Ext(event.Name) == ".deb" {
						_ = s.RebuildRepoMetadata(event.Name)
					}
					mutex.Unlock()
				}
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
		}
	}()
	return nil
}
