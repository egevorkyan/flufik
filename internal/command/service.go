package command

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/internal/repohttp"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/debrepository"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
)

type ServiceFlufikCommand struct {
	command *cobra.Command
}

func NewFlufikServiceCommand() *ServiceFlufikCommand {
	c := &ServiceFlufikCommand{
		command: &cobra.Command{
			Use:   "service",
			Short: "starts service",
		},
	}
	c.command.Run = c.Run
	return c
}

func (c *ServiceFlufikCommand) Run(command *cobra.Command, args []string) {
	if err := startService(); err != nil {
		logging.ErrorHandler("fatal: ", err)
	}

}

func startService() error {
	fmt.Println("Service Started!!!")
	serviceConfig := debrepository.NewServiceConfiguration()
	if err := crypto.GenerateKey(serviceConfig.PrivateKeyName, "", "", "", 0); err != nil {
		return err
	}

	if err := publishPublicKey(core.FlufikServiceWebHome(), serviceConfig.PrivateKeyName); err != nil {
		return err
	}
	if serviceConfig.EnableDirectoryWatching {
		if err := serviceConfig.DirectoryWatch(); err != nil {
			return err
		}
	}
	if err := serviceConfig.CreateDirectories(); err != nil {
		return err
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(core.FlufikServiceWebHome()))))
	http.Handle("/public", http.StripPrefix("/public", http.FileServer(http.Dir(core.FlufikServiceWebHome()))))
	http.Handle("/upload", repohttp.UploadHandler(serviceConfig))
	http.Handle("/delete", repohttp.DeleteHandler(serviceConfig))
	http.Handle("/unirepo", repohttp.GetApp())
	if serviceConfig.EnableSSL {
		if err := http.ListenAndServeTLS(":"+serviceConfig.ListenPort, filepath.Join(core.FlufikServiceConfigurationHome(), "server.crt"), filepath.Join(core.FlufikServiceConfigurationHome(), "server.key"), nil); err != nil {
			return err
		}
	} else {
		if err := http.ListenAndServe(":"+serviceConfig.ListenPort, nil); err != nil {
			return err
		}
	}
	return nil
}

func publishPublicKey(path string, keyName string) error {
	if err := os.MkdirAll(filepath.Join(path, "public"), 0755); err != nil {
		return err
	}
	if err := crypto.PublishPublicPGP(filepath.Join(path, "public"), keyName); err != nil {
		return err
	}
	return nil
}
