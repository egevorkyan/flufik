package crypto

import (
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/simpledb"
	"os"
	"path/filepath"
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
	DEFAULTTKNLEN  = 20
	DEFAULTTKNNUM  = 10
	DEFAULTTKNSYM  = 0
	DEFAULTTKNCAP  = 8
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

	//Saving private/public/passphrase keys in database
	if err := pgp.StoreKeysToDb(name, pgp.publicKey, pgp.privateKey, pgp.passPhrase); err != nil {
		return err
	}
	return nil
}

func CreateApiKey() error {
	token := PasswordGenerator(DEFAULTTKNLEN, DEFAULTTKNSYM, DEFAULTTKNNUM, DEFAULTTKNCAP)
	db := simpledb.NewSimpleDB()
	encodedToken := string(B64Encoder(token))
	if err := db.Insert("apikey", "", "", encodedToken); err != nil {
		return err
	}
	if err := SaveToFile(filepath.Join(core.FlufikServiceConfigurationHome(), "repo-token.txt"), B64Encoder(token)); err != nil {
		return err
	}
	db.CloseDb()
	return nil
}

// Gets Api key for authentication
func GetApiKey() (string, error) {
	db := simpledb.NewSimpleDB()
	var decodedKey string
	key, err := db.Get("apikey")
	if err != nil {
		return decodedKey, err
	}
	decodedByte, err := B64Decoder(string(key.TokenValue))
	if err != nil {
		return decodedKey, err
	}
	decodedKey = string(decodedByte)
	db.CloseDb()
	return decodedKey, nil
}
