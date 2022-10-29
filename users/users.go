package users

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/nosql"
	"os"
	"path/filepath"
	"strings"
)

const (
	USERCOLLECTION = "users"
	USERINDEXNAME  = "Username"
)

type Users struct{}

func NewUser() *Users {
	return &Users{}
}

func (u *Users) CreateUser(username string, mode string) error {
	data := make(map[string]interface{})
	tieDb, err := nosql.NewTieDot(USERCOLLECTION, USERINDEXNAME)
	if err != nil {
		return err
	}
	genQuery, err := tieDb.QueryGen(username, "eq", USERINDEXNAME)
	if err != nil {
		return err
	}
	docId, _, err := tieDb.Get(genQuery, USERCOLLECTION)
	if err != nil {
		return err
	}
	if docId == 0 {
		pwd := crypto.NewPwdGen(15, 3, 4, 3)
		pass, err := pwd.PasswordGenerator()
		if err != nil {
			return err
		}
		data["Username"] = username
		data["Password"] = pass
		data["Mode"] = mode
		err = tieDb.Insert(data, USERCOLLECTION)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) UpdateUser(username string) (string, error) {
	pwd := crypto.NewPwdGen(15, 3, 4, 3)
	pass, err := pwd.PasswordGenerator()
	if err != nil {
		return "", err
	}
	data := make(map[string]interface{})
	tieDb, err := nosql.NewTieDot(USERCOLLECTION, USERINDEXNAME)
	if err != nil {
		return "", err
	}
	genQuery, err := tieDb.QueryGen(username, "eq", USERINDEXNAME)
	if err != nil {
		return "", err
	}
	docId, value, err := tieDb.Get(genQuery, USERCOLLECTION)
	if err != nil {
		return "", err
	}
	if docId != 0 {
		data["Username"] = value["Username"]
		data["Password"] = pass
		data["Mode"] = value["Mode"]
		err = tieDb.Update(docId, data, USERCOLLECTION)
		if err != nil {
			return "", err
		}
	}
	return pass, nil
}

func (u *Users) Validate(username string, password string, mode string) (bool, error) {
	tieDb, err := nosql.NewTieDot(USERCOLLECTION, USERINDEXNAME)
	if err != nil {
		return false, err
	}
	genQuery, err := tieDb.QueryGen(username, "eq", USERINDEXNAME)
	if err != nil {
		return false, err
	}
	_, value, err := tieDb.Get(genQuery, USERCOLLECTION)
	if err != nil {
		return false, err
	}
	if strings.Compare(fmt.Sprint(value["Password"]), password) == 0 {
		if strings.Compare(fmt.Sprint(value["Mode"]), mode) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (u *Users) DumpUser(username string, fileName string) error {
	tieDb, err := nosql.NewTieDot(USERCOLLECTION, USERINDEXNAME)
	if err != nil {
		return err
	}
	genQuery, err := tieDb.QueryGen(username, "eq", USERINDEXNAME)
	if err != nil {
		return err
	}
	_, value, err := tieDb.Get(genQuery, USERCOLLECTION)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(core.FlufikServiceConfigurationHome(), fileName))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)
	_, err = fmt.Fprintf(f, "Username: %s Password: %s Mode: %s", value["Username"], value["Password"], value["Mode"])
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) DeleteUser(username string) error {
	tieDb, err := nosql.NewTieDot(USERCOLLECTION, USERINDEXNAME)
	if err != nil {
		return err
	}
	genQuery, err := tieDb.QueryGen(username, "eq", USERINDEXNAME)
	if err != nil {
		return err
	}
	docId, _, err := tieDb.Get(genQuery, USERCOLLECTION)
	if err != nil {
		return err
	}
	err = tieDb.Delete(docId, USERCOLLECTION)
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) GetUserPwd(username string) (string, error) {
	tieDb, err := nosql.NewTieDot(USERCOLLECTION, USERINDEXNAME)
	if err != nil {
		return "", err
	}
	genQuery, err := tieDb.QueryGen(username, "eq", USERINDEXNAME)
	if err != nil {
		return "", err
	}
	_, value, err := tieDb.Get(genQuery, USERCOLLECTION)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(value["Password"]), nil
}
