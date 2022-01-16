package crypto

import (
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"os"
)

const (
	DEFAULTKEYTYPE = "rsa"
	DEFAULTKEYBITS = 4096
	EXTENSION      = "pgp"
	PWDEXTENSION   = "txt"
	PRIVATEKEY     = "_priv"
	PUBLICKEY      = "_pub"
	PRIVATEKEYPWD  = "_pwd"
	DEFAULTPWDLEN  = 15
	DEFAULTPWDNUM  = 5
	DEFAULTPWDSYM  = 4
	DEFAULTPWDCAP  = 3
)

type FlufikPGP struct {
	privateKey string
	publicKey  string
	passPhrase string
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
	if err != nil {
		logging.ErrorHandler("armored private key failure ", err)
	}
	fpgp.privateKey = armoredPrivateKey
	armoredPublicKey, err := lockedKey.GetArmoredPublicKeyWithCustomHeaders(comment, version)
	fpgp.publicKey = armoredPublicKey
	fpgp.passPhrase = passphrase
	return &fpgp
}

func GenerateKey(name, email, comment, keyType string, bits int) error {
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

	//Generating automatically strong passPhase
	passPhrase := PasswordGenerator(DEFAULTPWDLEN, DEFAULTPWDSYM, DEFAULTPWDNUM, DEFAULTPWDCAP)

	pgp := NewPGP(name, email, comment, keyType, passPhrase, bits)

	//Commented section stores keys to files, soon will be decommissioned
	//privateKeyName := fmt.Sprintf("%s-%s", name, "priv")
	//publicKeyName := fmt.Sprintf("%s-%s", name, "pub")

	//if err := pgp.SaveKeys(privateKeyName, publicKeyName, EXTENSION); err != nil {
	//	return err
	//}

	//Saving private/public/passphrase keys in database
	if err := pgp.StoreKeysToDb(name); err != nil {
		return err
	}
	return nil
}

//This function will be completely decommissioned
/*func (f *FlufikPGP) SaveKeys(privName, pubName, ext string) error {
	privateKeyPath, publicKeyPath := core.FlufikKeyFileName(privName, pubName, ext)
	if err := ioutil.WriteFile(privateKeyPath, []byte(f.privateKey), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(publicKeyPath, []byte(f.publicKey), os.ModePerm); err != nil {
		return err
	}
	return nil
}*/
