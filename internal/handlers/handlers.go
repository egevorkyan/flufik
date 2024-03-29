package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/egevorkyan/flufik/pkg/config"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/egevorkyan/flufik/pkg/plugins/debrepository"
	"github.com/egevorkyan/flufik/pkg/plugins/rpmrepository"
	"github.com/egevorkyan/flufik/users"
	"github.com/gorilla/mux"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/basic"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"github.com/shaj13/go-guardian/store"
	"github.com/unrolled/render"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type deleteObject struct {
	Filename         string
	DistributionName string
	Arch             string
	Section          string
}

type Handler struct {
	cfg          *config.ServiceConfigBuilder
	deb          *debrepository.DebRepository
	yum          *rpmrepository.RpmRepo
	templatePath string
}

var authenticator auth.Authenticator
var cache store.Cache

func New(cfg *config.ServiceConfigBuilder, deb *debrepository.DebRepository, yum *rpmrepository.RpmRepo, templatePath string) *Handler {
	return &Handler{cfg: cfg, deb: deb, yum: yum, templatePath: templatePath}
}

func (h *Handler) UploadApt(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	archType := params["arch"][0]
	distroName := params["distro"][0]
	section := params["section"][0]
	r := render.New()
	if archType == "" {
		archType = "all"
	}
	if distroName == "" {
		distroName = "stable"
	}
	if section == "" {
		section = "main"
	}
	f, header, err := req.FormFile("file")
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	defer func(f multipart.File) {
		err = f.Close()
		if err != nil {
			_ = r.JSON(res, http.StatusInternalServerError, err.Error())
			return
		}
	}(f)
	path := filepath.Join(h.deb.ArchPath(distroName, section, archType), header.Filename)
	dest, err := os.Create(path)
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	defer func(dest *os.File) {
		err = dest.Close()
		if err != nil {
			_ = r.JSON(res, http.StatusInternalServerError, err.Error())
			return
		}
	}(dest)
	_, err = io.Copy(dest, f)
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	if err = h.deb.RebuildRepoMetadata(path); err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	_ = r.JSON(res, http.StatusOK, "package successfully uploaded")
}

func (h *Handler) UploadYum(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	arch := params["arch"][0]
	distroName := params["distro"][0]
	version := params["version"][0]
	r := render.New()

	f, header, err := req.FormFile("file")
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	defer func(f multipart.File) {
		err = f.Close()
		if err != nil {
			_ = r.JSON(res, http.StatusInternalServerError, err.Error())
			return
		}
	}(f)
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, f); err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	err = h.yum.Repository(buf.Bytes(), header.Filename, distroName, version, arch)
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	_ = r.JSON(res, http.StatusOK, "package successfully uploaded")
}

func (h *Handler) CreateUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	mode := vars["mode"]
	r := render.New()
	u := users.NewUser()
	err := u.CreateUser(username, mode)
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	uData, err := u.GetUserPwd(username)
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	_ = r.JSON(res, http.StatusOK, fmt.Sprintf("Password keep safe: %s", uData))
}

func (h *Handler) UpdateUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	r := render.New()
	u := users.NewUser()
	pwd, err := u.UpdateUser(username)
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	_ = r.JSON(res, http.StatusOK, fmt.Sprintf("Password updated keep safe: %s", pwd))
}

func (h *Handler) DeleteUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	r := render.New()
	u := users.NewUser()
	err := u.DeleteUser(username)
	if err != nil {
		_ = r.JSON(res, http.StatusInternalServerError, err.Error())
	}
	_ = r.JSON(res, http.StatusOK, fmt.Sprintf("User deleted: %s", username))
}

func (h *Handler) Middleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := authenticator.Authenticate(r)
		if err != nil {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		logger.InfoLog("User %s Authenticated\n", user.UserName())
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) SetupGoGuardian() {
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*10)

	basicStrategy := basic.New(validateUser, cache)
	tokenStrategy := bearer.New(bearer.NoOpAuthenticate, cache)

	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
	authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}

func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	u := users.NewUser()
	valid, err := u.Validate(userName, password, "admin")
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if valid {
		return auth.NewDefaultUser(userName, "1", nil, nil), nil
	}

	return nil, fmt.Errorf("invalid credentials")
}

func (h *Handler) GetHandler(f http.HandlerFunc) http.HandlerFunc {
	return f
}
