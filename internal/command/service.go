package command

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/internal/handlers"
	"github.com/egevorkyan/flufik/pkg/config"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/egevorkyan/flufik/pkg/plugins/debrepository"
	"github.com/egevorkyan/flufik/pkg/plugins/installer"
	"github.com/egevorkyan/flufik/pkg/plugins/rpmrepository"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
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
		logger.RaiseErr("fatal", err)
	}

}

func startService() error {
	cfg, err := config.GetServiceConfiguration()
	if err != nil {
		return err
	}
	deb := debrepository.NewServiceConfiguration(cfg)
	yum := rpmrepository.NewRpmBuilder(cfg)
	p := pgp.NewImportPGP()
	if err = os.MkdirAll(filepath.Join(core.FlufikServiceWebHome(cfg.RootRepoPath), "public"), 0755); err != nil {
		return fmt.Errorf("can not create directory: %v", err)
	}
	if err = p.PublishPublicPGP(filepath.Join(core.FlufikServiceWebHome(cfg.RootRepoPath), "public"), "flufik"); err != nil {
		return err
	}

	if err = getRepoInfo(core.FlufikServiceWebHome(cfg.RootRepoPath), cfg.RpmRepositoryName, cfg.PublicUrl); err != nil {
		return err
	}
	//if cfg.EnableDirectoryWatching {
	//	if err := deb.DirectoryWatch(); err != nil {
	//		return err
	//	}
	//}
	if err = deb.CreateDirectories(); err != nil {
		return err
	}
	if err = yum.CreateBaseDir(); err != nil {
		return err
	}
	handler := handlers.New(cfg, deb, yum, core.FlufikServiceWebHome(cfg.RootRepoPath))

	//Router start
	handler.SetupGoGuardian()
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(core.FlufikServiceWebHome(cfg.RootRepoPath))))).Methods("GET")
	router.HandleFunc("/upload/apt", handler.Middleware(handler.GetHandler(handler.UploadApt))).Methods("POST")
	router.HandleFunc("/upload/yum", handler.Middleware(handler.GetHandler(handler.UploadYum))).Methods("POST")
	router.HandleFunc("/user/add/{username}/{mode}", handler.Middleware(handler.GetHandler(handler.CreateUser))).Methods("POST")
	router.HandleFunc("/user/update/{username}", handler.Middleware(handler.GetHandler(handler.UpdateUser))).Methods("POST")
	router.HandleFunc("/user/delete/{username}", handler.Middleware(handler.GetHandler(handler.DeleteUser))).Methods("POST")
	chain := alice.New().Then(router)
	if err = http.ListenAndServe(":"+cfg.ListenPort, chain); err != nil {
		return err
	}
	logger.InfoLog("service started")
	return nil
}

func getRepoInfo(path string, repoName string, publicUrl string) error {
	i := installer.NewInstaller(repoName, publicUrl)
	if _, err := os.Stat(filepath.Join(path, "install")); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(path, "install"), 0755)
		if err != nil {
			return err
		}
	}
	err := i.GenerateShell(filepath.Join(path, "install"))
	if err != nil {
		return err
	}
	return nil
}
