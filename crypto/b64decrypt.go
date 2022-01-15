package crypto

import (
	"encoding/base64"
	"github.com/egevorkyan/flufik/core"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
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

func PasswordGenerator(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	pwdRune := []rune(password.String())
	rand.Shuffle(len(pwdRune), func(i, j int) {
		pwdRune[i], pwdRune[j] = pwdRune[j], pwdRune[i]
	})
	return string(pwdRune)
}
