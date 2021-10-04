package crypto

import (
	"fmt"
	"testing"
)

func TestFlufikReadPrivateKey(t *testing.T) {
	//read passphrased private key
	en, err := FlufikReadPrivateKey("flufik-priv.pgp", "Cisco2013")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(en)

}

func TestFlufikNewPGP(t *testing.T) {
	pgp := NewPGP("test", "test@outlook.com", "flufik generated pgp key", "rsa", "test123", 4096)
	fmt.Println(pgp.privateKey)
	fmt.Println(pgp.publicKey)

}
