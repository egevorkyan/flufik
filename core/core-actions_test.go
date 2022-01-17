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

func TestCheckArch(t *testing.T) {
	//Test #1: RPM package
	arch := CheckArch("flufik-0.2.7-1.el8.x86_64.rpm")
	fmt.Println(arch)

	//Test #2: DEB package
	arch = CheckArch("flufik_0.2.7-1_amd64.deb")
	fmt.Println(arch)
}
