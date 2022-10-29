package crypto

import (
	"github.com/sethvargo/go-password/password"
)

const (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
)

type pwdgen struct {
	passwordLength int
	minSpecialChar int
	minNum         int
	minUpperCase   int
}

func NewPwdGen(passwordLength, minSpecialChar, minNum, minUpperCase int) *pwdgen {
	return &pwdgen{
		passwordLength: passwordLength,
		minSpecialChar: minSpecialChar,
		minNum:         minNum,
		minUpperCase:   minUpperCase,
	}
}

func (p *pwdgen) PasswordGenerator() (string, error) {
	pwd, err := password.NewGenerator(&password.GeneratorInput{
		Symbols:      specialCharSet,
		Digits:       numberSet,
		LowerLetters: lowerCharSet,
		UpperLetters: upperCharSet,
	})
	if err != nil {
		return "", err
	}
	passPhrase, err := pwd.Generate(p.passwordLength, p.minNum, p.minSpecialChar, false, false)
	if err != nil {
		return "", err
	}
	return passPhrase, nil
}
