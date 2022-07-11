package crypto

import (
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/simpledb"
	"io/ioutil"
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
)

type FlufikPGP struct {
	privateKey string
	publicKey  string
	passPhrase string
}

// NewPGP - New PGP Builder
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

// GenerateKey - Generates new pgp key
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
	passPhrase, err := PasswordGenerator(DEFAULTPWDLEN, DEFAULTPWDSYM, DEFAULTPWDNUM, DEFAULTPWDCAP)
	if err != nil {
		return err
	}

	pgp := NewPGP(name, email, comment, keyType, passPhrase, bits)

	//Saving private/public/passphrase keys in database
	if err = pgp.StoreKeysToDb(name, pgp.publicKey, pgp.privateKey, pgp.passPhrase); err != nil {
		return err
	}
	return nil
}

// StoreKeysToDb - saves generated pgp keys to database
func (f *FlufikPGP) StoreKeysToDb(keyName string, publicKey string, privateKey string, pwd string) error {
	encodedPrivateKey := string(B64Encoder(privateKey))
	encodedPublicKey := string(B64Encoder(publicKey))
	encodedPwd := string(B64Encoder(pwd))
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	err = db.InsertPgpKeys(keyName, encodedPrivateKey, encodedPublicKey, encodedPwd)
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

// SavePgpKeyToFile - saves pgp key as files
func SavePgpKeyToFile(pgpKeyName string, location string) error {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	value, err := db.GetPgpByName(pgpKeyName)
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEY, EXTENSION)), []byte(value.PrivateKey)); err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PUBLICKEY, EXTENSION)), []byte(value.PublicKey)); err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEYPWD, PWDEXTENSION)), []byte(value.PassPhrase)); err != nil {
		return err
	}
	return nil
}

// PublishPublicPGP - publish pgp public key to be available via service url for Debian repositories
func PublishPublicPGP(filePath string, keyName string) error {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	publicKey, err := db.GetPgpByName(keyName)
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(filePath, fmt.Sprintf("%s%s.%s", keyName, PUBLICKEY, EXTENSION)), []byte(publicKey.PublicKey)); err != nil {
		return err
	}
	return nil
}

// SaveToFile - dumps to file
func SaveToFile(fileName string, encoded []byte) error {
	//decode to normal view
	decodedData, err := B64Decoder(string(encoded))
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(fileName, decodedData, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// ImportPgpKeys - Imports existing pgp key to flufik system
func ImportPgpKeys(name, private, public, passPhrase string) error {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	privateEncoded, err := readFile(private)
	if err != nil {
		return err
	}
	privateKeyEncoded := string(privateEncoded)
	publicEncoded, err := readFile(public)
	if err != nil {
		return err
	}
	publicKeyEncoded := string(publicEncoded)
	privatePwd := string(B64Encoder(passPhrase))
	err = db.InsertPgpKeys(name, privateKeyEncoded, publicKeyEncoded, privatePwd)
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

// readFile - reads pgp keys and returns value
func readFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	encoded := B64Encoder(string(data))
	return encoded, nil
}

// RemovePgpKeyFromDB - removes pgp key from database
func RemovePgpKeyFromDB(pgpName string) error {
	db, err := simpledb.OpenInternalDB(core.FlufikDbPath())
	if err != nil {
		return err
	}
	err = db.DeletePgpByName(pgpName)
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}
