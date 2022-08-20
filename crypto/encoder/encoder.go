package encoder

import (
	"encoding/base64"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"os"
	"path/filepath"
)

type encoder struct {
	logger    *logging.Logger
	debugging string
}

func NewEncoder(logger *logging.Logger, debugging string) *encoder {
	return &encoder{
		logger:    logger,
		debugging: debugging,
	}
}

func (e *encoder) B64Decoder(data string) ([]byte, error) {
	if e.debugging == "1" {
		e.logger.Info("Decoding base64 data")
	}
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func (e *encoder) B64Encoder(data string) []byte {
	if e.debugging == "1" {
		e.logger.Info("Encoding to base64 data")
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return []byte(encoded)
}

func (e *encoder) SaveB64DecodedData(data, fileName string) (string, error) {
	if e.debugging == "1" {
		e.logger.Info("Saving decoded data to file")
	}
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
