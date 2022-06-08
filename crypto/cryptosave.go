package crypto

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/plugins/simpledb"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (f *FlufikPGP) StoreKeysToDb(keyName string, publicKey string, privateKey string, pwd string) error {
	encodedPrivateKey := string(B64Encoder(privateKey))
	encodedPublicKey := string(B64Encoder(publicKey))
	encodedPwd := string(B64Encoder(pwd))
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
	if err := db.Insert(core.FLUFIKKEYDBTYPE, keyName, encodedPrivateKey, encodedPublicKey, encodedPwd); err != nil {
		return err
	}
	db.CloseDb()
	return nil
}

func SavePgpKeyToFile(pgpKeyName string, location string) error {
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
	value, err := db.GetKey(pgpKeyName)
	if err != nil {
		return err
	}
	db.CloseDb()
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEY, EXTENSION)), []byte(value.PrivateKeyValue)); err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PUBLICKEY, EXTENSION)), []byte(value.PublicKeyValue)); err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEYPWD, PWDEXTENSION)), []byte(value.TokenValue)); err != nil {
		return err
	}
	return nil
}

func PublishPublicPGP(filePath string, keyName string) error {
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
	publicKey, err := db.GetKey(keyName)
	db.CloseDb()
	if err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(filePath, fmt.Sprintf("%s%s.%s", keyName, PUBLICKEY, EXTENSION)), []byte(publicKey.PublicKeyValue)); err != nil {
		return err
	}
	return nil
}

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

func ImportPgpKeys(name, private, public, passPhrase string) error {
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
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
	if err = db.Insert(name, privateKeyEncoded, publicKeyEncoded, privatePwd); err != nil {
		return err
	}
	db.CloseDb()
	return nil
}
func readFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	encoded := B64Encoder(string(data))
	return encoded, nil
}

func RemovePgpKeyFromDB(pgpName string) error {
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
	if err := db.Delete(pgpName); err != nil {
		return err
	}
	db.CloseDb()
	return nil
}
