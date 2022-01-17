package crypto

import (
	"fmt"
	"testing"
)

func TestCheckSum(t *testing.T) {
	path := ""
	result, err := CheckSum(path)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result.Sha1)
	fmt.Println(result.Sha256)
	fmt.Println(result.Md5)
}
