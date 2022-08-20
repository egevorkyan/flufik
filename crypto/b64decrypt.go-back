package crypto

import (
	"encoding/base64"
	"github.com/egevorkyan/flufik/core"
	"github.com/sethvargo/go-password/password"
	"os"
	"path/filepath"
)

const (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
)

func B64Decoder(data string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func B64Encoder(data string) []byte {
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return []byte(encoded)
}

func SaveB64DecodedData(data, fileName string) (string, error) {
	decoded, err := B64Decoder(data)
	if err != nil {
		return "", err
	}
	absFileName := filepath.Join(core.FlufikServiceConfigurationHome(), fileName)
	f, err := os.Create(absFileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err = f.Write(decoded); err != nil {
		return "", err
	}
	if err = f.Sync(); err != nil {
		return "", err
	}
	return absFileName, nil
}

func PasswordGenerator(passwordLength, minSpecialChar, minNum, minUpperCase int) (string, error) {
	pwd, err := password.NewGenerator(&password.GeneratorInput{
		Symbols:      specialCharSet,
		Digits:       numberSet,
		LowerLetters: lowerCharSet,
		UpperLetters: upperCharSet,
	})
	if err != nil {
		return "", err
	}
	p, err := pwd.Generate(passwordLength, minNum, minSpecialChar, false, false)
	if err != nil {
		return "", err
	}
	return p, nil
}
