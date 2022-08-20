package pgp

import (
	"bytes"
	"crypto"
	"fmt"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/egevorkyan/flufik/crypto/encoder"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/nosql"
	"io"
	"os"
	"path/filepath"
	"unicode"
)

type signer struct {
	logger    *logging.Logger
	debugging string
}

func NewSigner(logger *logging.Logger, debugging string) *signer {
	return &signer{
		logger:    logger,
		debugging: debugging,
	}
}

func (s *signer) FlufikDebSigner(message io.Reader, privateKey string) ([]byte, error) {
	if s.debugging == "1" {
		s.logger.Info("debian package pgp key signer")
	}
	key, err := s.FlufikReadPrivateKey(privateKey)
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

func (s *signer) FlufikRpmSigner(privateKey string) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		if s.debugging == "1" {
			s.logger.Info("rpm package signer")
		}
		key, err := s.FlufikReadPrivateKey(privateKey)
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

func (s *signer) FlufikReadPrivateKey(privateKey string) (*openpgp.Entity, error) {
	if s.debugging == "1" {
		s.logger.Info("reading pgp private key")
	}
	tieDb, err := nosql.NewTieDot(PGPCOLLECTION, PGPINDEXNAME, s.logger, s.debugging)
	if err != nil {
		return nil, err
	}
	genQuery, err := tieDb.QueryGen(privateKey, "eq", PGPINDEXNAME)
	if err != nil {
		return nil, err
	}
	_, value, err := tieDb.Get(genQuery, PGPCOLLECTION)
	if err != nil {
		return nil, err
	}
	priv := fmt.Sprint(value["PrivateKey"])
	pwd := fmt.Sprint(value["PassPhrase"])
	e := encoder.NewEncoder(s.logger, s.debugging)
	decodedPrivateKey, err := e.B64Decoder(priv)
	if err != nil {
		return nil, err
	}
	passPhraseDecoded, err := e.B64Decoder(pwd)
	if err != nil {
		return nil, err
	}

	var entityList openpgp.EntityList

	if s.FlufikCheckPGPKeyType(decodedPrivateKey) {
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

func (s *signer) FlufikDecryptPrivateKey(keyName, passPhrase string) (*openpgp.Entity, error) {
	if s.debugging == "1" {
		s.logger.Info("decrypt private key")
	}
	e := encoder.NewEncoder(s.logger, s.debugging)
	tieDb, err := nosql.NewTieDot(PGPCOLLECTION, PGPINDEXNAME, s.logger, s.debugging)
	if err != nil {
		return nil, err
	}
	genQuery, err := tieDb.QueryGen(keyName, "eq", PGPINDEXNAME)
	if err != nil {
		return nil, err
	}
	_, value, err := tieDb.Get(genQuery, PGPCOLLECTION)
	if err != nil {
		return nil, err
	}
	privateKey, err := e.B64Decoder(fmt.Sprint(value["PrivateKey"]))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	var entityList openpgp.EntityList

	if s.FlufikCheckPGPKeyType(privateKey) {
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

func (s *signer) FlufikCheckPGPKeyType(pgpKey []byte) bool {
	if s.debugging == "1" {
		s.logger.Info("identifying if pgp key is encrypted or not")
	}
	for pgp := 0; pgp < len(pgpKey); pgp++ {
		if pgpKey[pgp] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func (s *signer) SignRelease(privateKey string, fileName string) error {
	if s.debugging == "1" {
		s.logger.Info("signing release files for debian repository")
	}
	key, err := s.FlufikReadPrivateKey(privateKey)
	if err != nil {
		return err
	}
	workingDirectory := filepath.Dir(fileName)
	releaseFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening release file (%s) for writing: %s", fileName, err)
	}

	releaseGpg, err := os.Create(filepath.Join(workingDirectory, "Release.gpg"))
	if err != nil {
		return fmt.Errorf("error creating Release.pgp file for writing: %s", err)
	}
	defer releaseGpg.Close()
	if err = openpgp.ArmoredDetachSign(releaseGpg, key, releaseFile, &packet.Config{
		DefaultHash: crypto.SHA256,
	}); err != nil {
		return fmt.Errorf("armored detached sign failure: %s", err)
	}
	releaseFile.Seek(0, 0)
	inlineRelease, err := os.Create(filepath.Join(workingDirectory, "InRelease"))
	if err != nil {
		return fmt.Errorf("error creating InRelease file for writing: %s", err)
	}
	defer inlineRelease.Close()

	writer, err := clearsign.Encode(inlineRelease, key.PrivateKey, nil)
	if err != nil {
		return fmt.Errorf("error signing InRelease file : %s", err)
	}
	io.Copy(writer, releaseFile)
	writer.Close()
	return nil
}
