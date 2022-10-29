package encoder

import (
	"encoding/base64"
	"github.com/egevorkyan/flufik/core"
	"os"
	"path/filepath"
)

type encoder struct{}

func NewEncoder() *encoder {
	return &encoder{}
}

func (e *encoder) B64Decoder(data string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func (e *encoder) B64Encoder(data string) []byte {
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return []byte(encoded)
}

func (e *encoder) SaveB64DecodedData(data, fileName string) (string, error) {
	decoded, err := e.B64Decoder(data)
	if err != nil {
		return "", err
	}
	absFileName := filepath.Join(core.FlufikServiceConfigurationHome(), fileName)
	f, err := os.Create(absFileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err = f.Write(decoded); err != nil {
		return "", err
	}
	if err = f.Sync(); err != nil {
		return "", err
	}
	return absFileName, nil
}
