package repohttp

import (
	"encoding/json"
	"fmt"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/plugins/debrepository"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type deleteObject struct {
	Filename         string
	DistributionName string
	Arch             string
	Section          string
}

func UploadHandler(serviceConfig *debrepository.ServiceConfigBuilder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
		if serviceConfig.EnableAPIKeys {
			apiKey := r.URL.Query().Get("key")
			if apiKey == "" {
				http.Error(w, "api key not present", http.StatusUnauthorized)
				return
			}
			if !validateAPIkey(apiKey) {
				http.Error(w, "api key not valid", http.StatusUnauthorized)
				return
			}
		}

		archType := r.URL.Query().Get("arch")
		if archType == "" {
			archType = "all"
		}
		distroName := r.URL.Query().Get("distro")
		if distroName == "" {
			distroName = "stable"
		}
		section := r.URL.Query().Get("section")
		if section == "" {
			section = "main"
		}
		reader, err := r.MultipartReader()
		if err != nil {
			httpErrorFormat(w, "error creating multipart reader: %s", err)
			return
		}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if part.FileName() == "" {
				continue
			}

			path := filepath.Join(serviceConfig.ArchPath(distroName, section, archType), part.FileName())
			fmt.Println(path)
			dst, err := os.Create(path)

			if err != nil {
				httpErrorFormat(w, "error creating deb file: %s", err)
				return
			}
			if _, err := io.Copy(dst, part); err != nil {
				httpErrorFormat(w, "error writing deb file: %s", err)
				return
			}

			dst.Close()
			if !serviceConfig.EnableDirectoryWatching {
				if err := serviceConfig.RebuildRepoMetadata(path); err != nil {
					httpErrorFormat(w, "error rebuild repository: %s", err)
				}
			}

		}
		w.WriteHeader(http.StatusOK)
	})
}

func DeleteHandler(serviceConfig *debrepository.ServiceConfigBuilder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		if r.Method != "DELETE" {
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
		if serviceConfig.EnableAPIKeys {
			apiKey := r.URL.Query().Get("key")
			if apiKey == "" {
				http.Error(w, "api key not present", http.StatusUnauthorized)
				return
			}
			if !validateAPIkey(apiKey) {
				http.Error(w, "api key not valid", http.StatusUnauthorized)
				return
			}
		}
		var toDelete deleteObject
		if err := json.NewDecoder(r.Body).Decode(&toDelete); err != nil {
			httpErrorFormat(w, "failed to decode json: %s", err)
			return
		}
		debPath := filepath.Join(serviceConfig.ArchPath(toDelete.DistributionName, toDelete.Section, toDelete.Arch), toDelete.Filename)
		if err := os.Remove(debPath); err != nil {
			httpErrorFormat(w, "failed to delete: %s", err)
			return
		}

	})
}

//Universal applications' repository handlers
func GetApp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
		_, err := io.WriteString(w, runtime.GOOS)
		if err != nil {
			httpErrorFormat(w, "failed to show: %s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func httpErrorFormat(w http.ResponseWriter, format string, a ...interface{}) {
	err := fmt.Errorf(format, a...)
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func validateAPIkey(key string) bool {
	k, _ := crypto.GetApiKey()
	if k == key {
		return true
	} else {
		return false
	}
}
