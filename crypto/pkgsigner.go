package crypto

import (
	"bytes"
	"crypto"
	"fmt"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/plugins/simpledb"
	"io"
	"unicode"
)

// FlufikDebSigner - Debian package signer
func FlufikDebSigner(message io.Reader, privateKey string) ([]byte, error) {
	key, err := FlufikReadPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	var signature bytes.Buffer

	if err = openpgp.ArmoredDetachSign(&signature, key, message, &packet.Config{
		DefaultHash: crypto.SHA256,
	}); err != nil {
		return nil, fmt.Errorf("armored detached sign failure: %w", err)
	}

	return signature.Bytes(), nil
}

// FlufikRpmSigner - rpm package signer
func FlufikRpmSigner(privateKey string) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		key, err := FlufikReadPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		var signature bytes.Buffer
		if err = openpgp.DetachSign(&signature, key, bytes.NewReader(data), &packet.Config{
			DefaultHash: crypto.SHA256,
		}); err != nil {
			return nil, err
		}
		return signature.Bytes(), nil
	}
}

func FlufikReadPrivateKey(privateKey string) (*openpgp.Entity, error) {
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
	privateEncoded, err := db.GetKey(privateKey)
	if err != nil {
		return nil, err
	}
	db.CloseDb()
	priv := privateEncoded.PrivateKeyValue
	pwd := privateEncoded.TokenValue
	decodedPrivateKey, err := B64Decoder(priv)
	if err != nil {
		return nil, err
	}
	passPhraseDecoded, err := B64Decoder(pwd)
	if err != nil {
		return nil, err
	}

	var entityList openpgp.EntityList

	if FlufikCheckPGPKeyType(decodedPrivateKey) {
		entityList, err = openpgp.ReadArmoredKeyRing(bytes.NewReader(decodedPrivateKey))
		if err != nil {
			return nil, fmt.Errorf("decoding armored PGP keyring failure %w", err)
		}
	} else {
		entityList, err = openpgp.ReadKeyRing(bytes.NewReader(decodedPrivateKey))
		if err != nil {
			return nil, fmt.Errorf("decoding failure %w", err)
		}
	}

	key := entityList[0]

	if key.PrivateKey == nil {
		return nil, fmt.Errorf("no private key")
	}

	if key.PrivateKey.Encrypted {
		if string(passPhraseDecoded) == "" {
			return nil, fmt.Errorf("key encrypted, passphrase not provided")
		}
		if err = key.PrivateKey.Decrypt(passPhraseDecoded); err != nil {
			return nil, fmt.Errorf("failure decrypting private key: %w", err)
		}
		for _, subKey := range key.Subkeys {
			if subKey.PrivateKey != nil {
				if err = subKey.PrivateKey.Decrypt(passPhraseDecoded); err != nil {
					return nil, fmt.Errorf("failure decrypting sub private key: %w", err)
				}
			}
		}

	}
	return key, nil
}

func FlufikDecryptPrivateKey(keyName, passPhrase, dbName string) (*openpgp.Entity, error) {
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
	encodedPrivateKey, err := db.GetKey(keyName)
	if err != nil {
		return nil, fmt.Errorf("fatal: %w", err)
	}
	db.CloseDb()

	privateKey, err := B64Decoder(encodedPrivateKey.PrivateKeyValue)
	if err != nil {
		return nil, fmt.Errorf("fatal: %w", err)
	}

	var entityList openpgp.EntityList

	if FlufikCheckPGPKeyType(privateKey) {
		entityList, err = openpgp.ReadArmoredKeyRing(bytes.NewReader(privateKey))
		if err != nil {
			return nil, fmt.Errorf("decoding armored PGP keyring failure %w", err)
		}
	} else {
		entityList, err = openpgp.ReadKeyRing(bytes.NewReader(privateKey))
		if err != nil {
			return nil, fmt.Errorf("decoding failure %w", err)
		}
	}

	key := entityList[0]

	if key.PrivateKey == nil {
		return nil, fmt.Errorf("no private key")
	}

	if key.PrivateKey.Encrypted {
		if passPhrase == "" {
			return nil, fmt.Errorf("key encrypted, passphrase not provided")
		}

		pwd := []byte(passPhrase)
		if err = key.PrivateKey.Decrypt(pwd); err != nil {
			return nil, fmt.Errorf("failure decrypting private key: %w", err)
		}
		for _, subKey := range key.Subkeys {
			if subKey.PrivateKey != nil {
				if err = subKey.PrivateKey.Decrypt(pwd); err != nil {
					return nil, fmt.Errorf("failure decrypting sub private key: %w", err)
				}
			}
		}

	}

	return key, nil
}

func FlufikCheckPGPKeyType(pgpKey []byte) bool {
	for pgp := 0; pgp < len(pgpKey); pgp++ {
		if pgpKey[pgp] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
