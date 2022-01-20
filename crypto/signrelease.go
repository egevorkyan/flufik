package crypto

import (
	"crypto"
	"fmt"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"io"
	"os"
	"path/filepath"
)

//Signs Release.pgp and InRelease
func SignRelease(privateKey string, fileName string) error {
	key, err := FlufikReadPrivateKey(privateKey)
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
