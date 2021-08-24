package crypto

import (
	"bytes"
	"crypto"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"io/ioutil"
)

// FlufikDebSigner - Debian package signer
func FlufikDebSigner(message io.Reader, keyFile string, passPhrase string) ([]byte, error) {
	key, err := flufikReadPrivateKey(keyFile, passPhrase)
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
func FlufikRpmSigner(keyFile string, passPhrase string) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		key, err := flufikReadPrivateKey(keyFile, passPhrase)
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

func flufikReadPrivateKey(keyFile string, passPhrase string) (*openpgp.Entity, error) {
	privateKeyFile, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("reading PGP private key failure %w", err)
	}

	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(privateKeyFile))
	if err != nil {
		return nil, fmt.Errorf("decoding armored PGP keyring failure %w", err)
	}

	key := entityList[0]

	if key.PrivateKey == nil {
		return nil, fmt.Errorf("no private key")
	}

	if key.PrivateKey.Encrypted {
		if len(passPhrase) == 0 {
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
