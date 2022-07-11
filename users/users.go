package users

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/plugins/simpledb"
	"os"
	"path/filepath"
	"strings"
)

type Users struct{}

func NewUser() *Users {
	return &Users{}
}

func (u *Users) CreateUser(username string, mode string) error {
	pwd, err := crypto.PasswordGenerator(15, 3, 4, 3)
	if err != nil {
		return err
	}
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	err = db.InsertUsers(username, pwd, mode)
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) UpdateUser(username string) (string, error) {
	pwd, err := crypto.PasswordGenerator(15, 3, 4, 3)
	if err != nil {
		return "", err
	}
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return "", err
	}
	err = db.UpdateUserByName(username, pwd)
	if err != nil {
		return "", err
	}
	err = db.Close()
	if err != nil {
		return "", err
	}
	return pwd, nil
}

func (u *Users) Validate(username string, password string, mode string) (bool, error) {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return false, err
	}
	userData, err := db.GetUserByName(username)
	if err != nil {
		return false, err
	}
	if strings.Compare(userData.Password, password) == 0 {
		if strings.Compare(userData.Mode, mode) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (u *Users) DumpUser(username string, fileName string) error {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	userData, err := db.GetUserByName(username)
	if err != nil {
		return err
	}
	err = db.Close()
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
	_, err = fmt.Fprintf(f, "Username: %s Password: %s Mode: %s", userData.UserName, userData.Password, userData.Mode)
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) DeleteUser(username string) error {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	err = db.DeleteUserByName(username)
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) GetUserPwd(username string) (string, error) {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return "", err
	}
	userData, err := db.GetUserByName(username)
	if err != nil {
		return "", err
	}
	err = db.Close()
	if err != nil {
		return "", err
	}
	return userData.Password, nil
}
