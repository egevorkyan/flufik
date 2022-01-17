package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

type FlufikChecksum struct {
	Sha1   string
	Sha256 string
	Md5    string
}

func CheckSum(file string) (hash FlufikChecksum, err error) {
	sha256Hash := sha256.New()
	sha1Hash := sha1.New()
	md5Hash := md5.New()

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return FlufikChecksum{}, fmt.Errorf("can not open file: %w", err)
	}

	sha1Hash.Write(f)
	sha256Hash.Write(f)
	md5Hash.Write(f)

	hash.Sha1 = hex.EncodeToString(sha1Hash.Sum(nil))
	hash.Sha256 = hex.EncodeToString(sha256Hash.Sum(nil))
	hash.Md5 = hex.EncodeToString(md5Hash.Sum(nil))

	return hash, nil
}
