package pgp

import (
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/egevorkyan/flufik/core"
	plugin "github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/crypto/encoder"
	"github.com/egevorkyan/flufik/pkg/nosql"
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
	PGPCOLLECTION  = "pgp"
	PGPINDEXNAME   = "KeyName"
)

type pgpKey struct {
	name    string
	email   string
	comment string
	keyType string
	bits    int
}

func NewPGP(name, email, comment, keyType string, bits int) *pgpKey {
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
	return &pgpKey{
		name:    name,
		email:   email,
		comment: comment,
		keyType: keyType,
		bits:    bits,
	}
}

func NewImportPGP() *pgpKey {
	return &pgpKey{}
}

func (p *pgpKey) GeneratePgpKey() error {
	pwd := plugin.NewPwdGen(15, 3, 4, 3)

	pass, err := pwd.PasswordGenerator()
	if err != nil {
		return err
	}
	key, err := crypto.GenerateKey(p.name, p.email, p.keyType, p.bits)
	if err != nil {
		return fmt.Errorf("pgp generate failure %v", err)
	}

	lockedKey, err := key.Lock([]byte(pass))
	if err != nil {
		return fmt.Errorf("pgp lock failure %v", err)
	}

	version := fmt.Sprintf("flufik-%s", core.Version)

	armoredPrivateKey, err := lockedKey.ArmorWithCustomHeaders(p.comment, version)
	if err != nil {
		return fmt.Errorf("armored private key failure %v", err)
	}
	armoredPublicKey, err := lockedKey.GetArmoredPublicKeyWithCustomHeaders(p.comment, version)
	//Saving private/public/passphrase keys in database
	if err = p.StoreKeysToDb(p.name, armoredPublicKey, armoredPrivateKey, pass); err != nil {
		return err
	}
	return nil
}

func (p *pgpKey) StoreKeysToDb(keyName string, publicKey string, privateKey string, pwd string) error {
	data := make(map[string]interface{})
	e := encoder.NewEncoder()
	encodedPrivateKey := string(e.B64Encoder(privateKey))
	encodedPublicKey := string(e.B64Encoder(publicKey))
	encodedPwd := string(e.B64Encoder(pwd))
	data["KeyName"] = keyName
	data["PrivateKey"] = encodedPrivateKey
	data["PublicKey"] = encodedPublicKey
	data["PassPhrase"] = encodedPwd
	tieDb, err := nosql.NewTieDot(PGPCOLLECTION, PGPINDEXNAME)
	if err != nil {
		return err
	}
	err = tieDb.Insert(data, PGPCOLLECTION)
	if err != nil {
		return err
	}
	return nil
}

func (p *pgpKey) SavePgpKeyToFile(pgpKeyName string, location string) error {
	tieDb, err := nosql.NewTieDot(PGPCOLLECTION, PGPINDEXNAME)
	if err != nil {
		return err
	}
	genQuery, err := tieDb.QueryGen(pgpKeyName, "eq", PGPINDEXNAME)
	if err != nil {
		return err
	}
	_, value, err := tieDb.Get(genQuery, PGPCOLLECTION)
	if err != nil {
		return err
	}
	if err = p.saveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEY, EXTENSION)), []byte(fmt.Sprint(value["PrivateKey"]))); err != nil {
		return err
	}
	if err = p.saveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PUBLICKEY, EXTENSION)), []byte(fmt.Sprint(value["PublicKey"]))); err != nil {
		return err
	}
	if err = p.saveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEYPWD, PWDEXTENSION)), []byte(fmt.Sprint(value["PassPhrase"]))); err != nil {
		return err
	}
	return nil
}

func (p *pgpKey) ImportPgpKeys(name, private, public, passPhrase string) error {
	data := make(map[string]interface{})
	e := encoder.NewEncoder()

	privateEncoded, err := p.readFile(private)
	if err != nil {
		return err
	}
	privateKeyEncoded := string(privateEncoded)
	publicEncoded, err := p.readFile(public)
	if err != nil {
		return err
	}
	publicKeyEncoded := string(publicEncoded)
	privatePwd := string(e.B64Encoder(passPhrase))
	data["KeyName"] = name
	data["PrivateKey"] = privateKeyEncoded
	data["PublicKey"] = publicKeyEncoded
	data["PassPhrase"] = privatePwd
	tieDb, err := nosql.NewTieDot(PGPCOLLECTION, PGPINDEXNAME)
	if err != nil {
		return err
	}
	genQuery, err := tieDb.QueryGen(name, "eq", PGPINDEXNAME)
	if err != nil {
		return err
	}
	docId, _, err := tieDb.Get(genQuery, PGPCOLLECTION)
	if err != nil {
		return err
	}
	if docId == 0 {
		err = tieDb.Insert(data, PGPCOLLECTION)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *pgpKey) RemovePgpKeyFromDB(keyName string) error {
	tieDb, err := nosql.NewTieDot(PGPCOLLECTION, PGPINDEXNAME)
	if err != nil {
		return err
	}
	genQuery, err := tieDb.QueryGen(keyName, "eq", PGPINDEXNAME)
	if err != nil {
		return err
	}
	docId, _, err := tieDb.Get(genQuery, PGPCOLLECTION)
	if err != nil {
		return err
	}
	err = tieDb.Delete(docId, PGPCOLLECTION)
	if err != nil {
		return err
	}
	return nil
}

func (p *pgpKey) PublishPublicPGP(filePath string, keyName string) error {
	tieDb, err := nosql.NewTieDot(PGPCOLLECTION, PGPINDEXNAME)
	if err != nil {
		return err
	}
	genQuery, err := tieDb.QueryGen(keyName, "eq", PGPINDEXNAME)
	if err != nil {
		return err
	}
	_, value, err := tieDb.Get(genQuery, PGPCOLLECTION)
	if err != nil {
		return err
	}
	if err = p.saveToFile(filepath.Join(filePath, fmt.Sprintf("%s%s.%s", keyName, PUBLICKEY, EXTENSION)), []byte(fmt.Sprint(value["PublicKey"]))); err != nil {
		return err
	}
	return nil
}

func (p *pgpKey) saveToFile(fileName string, encoded []byte) error {
	e := encoder.NewEncoder()
	//decode to normal view
	decodedData, err := e.B64Decoder(string(encoded))
	if err != nil {
		return err
	}
	if err = os.WriteFile(fileName, decodedData, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (p *pgpKey) readFile(path string) ([]byte, error) {
	e := encoder.NewEncoder()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	encoded := e.B64Encoder(string(data))
	return encoded, nil
}
