package command

import (
	"bytes"
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/internal/handlers"
	"github.com/egevorkyan/flufik/pkg/config"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/debrepository"
	"github.com/egevorkyan/flufik/pkg/plugins/rpmrepository"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/spf13/cobra"
	"html/template"
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
	cfg := config.GetServiceConfiguration()
	deb := debrepository.NewServiceConfiguration(cfg)
	yum := rpmrepository.NewRpmBuilder(cfg)
	if err := crypto.GenerateKey(cfg.PrivateKeyName, "", "", "", 0); err != nil {
		return err
	}

	if err := publishPublicKey(core.FlufikServiceWebHome(cfg.RootRepoPath), cfg.PrivateKeyName); err != nil {
		return err
	}
	if err := getRepoInfo(core.FlufikServiceWebHome(cfg.RootRepoPath), cfg.RpmRepositoryName); err != nil {
		return err
	}
	if cfg.EnableDirectoryWatching {
		if err := deb.DirectoryWatch(); err != nil {
			return err
		}
	}
	if err := deb.CreateDirectories(); err != nil {
		return err
	}
	if err := yum.CreateBaseDir(); err != nil {
		return err
	}
	handler := handlers.New(cfg, deb, yum)

	//Router start
	handler.SetupGoGuardian()
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(core.FlufikServiceWebHome(cfg.RootRepoPath))))).Methods("GET")
	router.HandleFunc("/upload", handler.Middleware(handler.GetHandler(handler.Upload))).Methods("POST")
	router.HandleFunc("/user/add/{username}/{mode}", handler.Middleware(handler.GetHandler(handler.CreateUser))).Methods("POST")
	router.HandleFunc("/user/update/{username}", handler.Middleware(handler.GetHandler(handler.UpdateUser))).Methods("POST")
	router.HandleFunc("/user/delete/{username}", handler.Middleware(handler.GetHandler(handler.DeleteUser))).Methods("POST")
	chain := alice.New().Then(router)
	if cfg.EnableSSL {
		if err := http.ListenAndServeTLS(":"+cfg.ListenPort, filepath.Join(core.FlufikServiceConfigurationHome(), "server.crt"), filepath.Join(core.FlufikServiceConfigurationHome(), "server.key"), chain); err != nil {
			return err
		}
	} else {
		if err := http.ListenAndServe(":"+cfg.ListenPort, chain); err != nil {
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

func getRepoInfo(path string, repoName string) error {
	proto := "http(s)"
	data := struct {
		Reponame string
		URL      string
	}{
		Reponame: repoName,
		URL:      proto + "://fdqn/" + repoName + "/",
	}
	t := `[{{.Reponame}}]
name={{.Reponame}}
baseurl={{.URL}}
enabled=1
gpgcheck=0
priority=1`
	templ, err := template.New("repofile").Parse(t)
	if err != nil {
		return err
	}
	var n bytes.Buffer
	err = templ.Execute(&n, data)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(path, "public", "howto.txt"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	dRepo := fmt.Sprintf("Debian Based Repository configuration:\nExample: deb %s://fqdn/ stable main\n\nRedHat Based Repository configuration:\n%s", proto, n.String())

	_, err = fmt.Fprint(f, dRepo)
	if err != nil {
		return err
	}
	return nil
}
