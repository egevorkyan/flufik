package crypto

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/plugins/badgerdb"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (f *FlufikPGP) StoreKeysToDb(keyName string) error {
	var kvdb = make(map[string]string)
	kvdb[fmt.Sprintf("%s%s", keyName, PRIVATEKEY)] = string(B64Encoder(f.privateKey))
	kvdb[fmt.Sprintf("%s%s", keyName, PUBLICKEY)] = string(B64Encoder(f.publicKey))
	kvdb[fmt.Sprintf("%s%s", keyName, PRIVATEKEYPWD)] = string(B64Encoder(f.passPhrase))
	db := badgerdb.NewFlufikBadgerDB(core.FlufikKeyDbPath())
	if err := db.UpdateDb(kvdb); err != nil {
		return err
	}
	db.Close()
	return nil
}

func SavePgpKeyToFile(pgpKeyName string, location string) error {
	db := badgerdb.NewFlufikBadgerDB(core.FlufikKeyDbPath())
	privateKey, err := db.Get(fmt.Sprintf("%s%s", pgpKeyName, PRIVATEKEY))
	if err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEY, EXTENSION)), privateKey); err != nil {
		return err
	}
	publicKey, err := db.Get(fmt.Sprintf("%s%s", pgpKeyName, PUBLICKEY))
	if err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PUBLICKEY, EXTENSION)), publicKey); err != nil {
		return err
	}
	passPhrase, err := db.Get(fmt.Sprintf("%s%s", pgpKeyName, PRIVATEKEYPWD))
	if err != nil {
		return err
	}
	if err = SaveToFile(filepath.Join(location, fmt.Sprintf("%s%s.%s", pgpKeyName, PRIVATEKEYPWD, PWDEXTENSION)), passPhrase); err != nil {
		return err
	}
	db.Close()
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
	var kvdb = make(map[string]string)
	db := badgerdb.NewFlufikBadgerDB(core.FlufikKeyDbPath())
	privateEncoded, err := readFile(private)
	if err != nil {
		return err
	}
	kvdb[fmt.Sprintf("%s%s", name, PRIVATEKEY)] = string(privateEncoded)
	publicEncoded, err := readFile(public)
	if err != nil {
		return err
	}
	kvdb[fmt.Sprintf("%s%s", name, PUBLICKEY)] = string(publicEncoded)
	kvdb[fmt.Sprintf("%s%s", name, PRIVATEKEYPWD)] = string(B64Encoder(passPhrase))
	if err = db.UpdateDb(kvdb); err != nil {
		return err
	}
	db.Close()
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
	db := badgerdb.NewFlufikBadgerDB(core.FlufikKeyDbPath())
	privateKey := fmt.Sprintf("%s%s", pgpName, PRIVATEKEY)
	publicKey := fmt.Sprintf("%s%s", pgpName, PUBLICKEY)
	pwdPrivate := fmt.Sprintf("%s%s", pgpName, PRIVATEKEYPWD)
	_, err := db.Get(privateKey)
	if err != nil {
		return err
	} else {
		if err = db.Remove(privateKey); err != nil {
			return err
		}
	}
	_, err = db.Get(publicKey)
	if err != nil {
		return err
	} else {
		if err := db.Remove(publicKey); err != nil {
			return err
		}
	}
	_, err = db.Get(pwdPrivate)
	if err != nil {
		return err
	} else {
		if err := db.Remove(pwdPrivate); err != nil {
			return err
		}
	}
	db.Close()
	return nil
}