package core

import (
	"bytes"
	"crypto"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"io/ioutil"
)

func PGPArmoredSign(message io.Reader, keyFile string) ([]byte, error) {
	key, err := readPrivateKey(keyFile)
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

func readPrivateKey(keyFile string) (*openpgp.Entity, error) {
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
		return nil, fmt.Errorf("key encrypted, passphrase required")
	}

	return key, nil
}
