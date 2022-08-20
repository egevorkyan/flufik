package crypto

import (
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/sethvargo/go-password/password"
)

const (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
)

type pwdgen struct {
	logger         *logging.Logger
	debugging      string
	passwordLength int
	minSpecialChar int
	minNum         int
	minUpperCase   int
}

func NewPwdGen(passwordLength, minSpecialChar, minNum, minUpperCase int, logger *logging.Logger, debugging string) *pwdgen {
	return &pwdgen{
		passwordLength: passwordLength,
		minSpecialChar: minSpecialChar,
		minNum:         minNum,
		minUpperCase:   minUpperCase,
		logger:         logger,
		debugging:      debugging,
	}
}

func (p *pwdgen) PasswordGenerator() (string, error) {
	if p.debugging == "1" {
		p.logger.Info("generating password")
	}
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
