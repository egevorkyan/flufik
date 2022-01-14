package crypto

import (
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"io/ioutil"
	"os"
)

const (
	DEFAULTKEYTYPE = "rsa"
	DEFAULTKEYBITS = 4096
	EXTENSION      = "pgp"
)

type FlufikPGP struct {
	privateKey string
	publicKey  string
}

func NewPGP(name, email, comment, keyType, passphrase string, bits int) *FlufikPGP {
	var fpgp FlufikPGP
	key, err := crypto.GenerateKey(name, email, keyType, bits)
	if err != nil {
		logging.ErrorHandler("pgp generate failure ", err)
	}

	lockedKey, err := key.Lock([]byte(passphrase))
	if err != nil {
		logging.ErrorHandler("pgp lock failure ", err)
	}

	version := fmt.Sprintf("flufik-%s", core.Version)

	armoredPrivateKey, err := lockedKey.ArmorWithCustomHeaders(comment, version)
	//armoredPrivateKey, err := lockedKey.Armor()
	if err != nil {
		logging.ErrorHandler("armored private key failure ", err)
	}
	fpgp.privateKey = armoredPrivateKey
	armoredPublicKey, err := lockedKey.GetArmoredPublicKeyWithCustomHeaders(comment, version)
	//armoredPublicKey, err := lockedKey.GetArmoredPublicKey()
	fpgp.publicKey = armoredPublicKey
	return &fpgp
}

func (f *FlufikPGP) SaveKeys(privName, pubName, ext string) error {
	privateKeyPath, publicKeyPath := core.FlufikKeyFileName(privName, pubName, ext)
	if err := ioutil.WriteFile(privateKeyPath, []byte(f.privateKey), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(publicKeyPath, []byte(f.publicKey), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func GenerateKey(name, email, comment, keyType, passphrase string, bits int) error {
	if name == "" {
		h, _ := os.Hostname()
		name = fmt.Sprintf("flufik-%s", h)
	}
	if email == "" {
		email = "flufik@flufik.com"
	}
	if keyType == "" {
		keyType = DEFAULTKEYTYPE
	}
	if bits == 0 {
		bits = DEFAULTKEYBITS
	}
	if comment == "" {
		comment = fmt.Sprintf("Flufik - Type: %s %v Bits", keyType, bits)
	}
	if passphrase == "" {
		return fmt.Errorf("Passphrase is empty, please provide passphrase")
	}
	pgp := NewPGP(name, email, comment, keyType, passphrase, bits)

	privateKeyName := fmt.Sprintf("%s-%s", name, "priv")
	publicKeyName := fmt.Sprintf("%s-%s", name, "pub")

	if err := pgp.SaveKeys(privateKeyName, publicKeyName, EXTENSION); err != nil {
		return err
	}
	return nil
}
