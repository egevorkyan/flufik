package core

import (
	"fmt"
	"testing"
)

func TestOpenFile(t *testing.T) {
	path := ""
	fileName := "parsed.txt"
	f, err := OpenFile(fileName, path)
	if err != nil {
		t.Errorf("failed %w", err)
	}
	fmt.Println(f.Name())
}

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
