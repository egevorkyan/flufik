package rpmrepository

import (
	"fmt"
	"github.com/egevorkyan/flufik/pkg/config"
	"io/ioutil"
	"testing"
)

func TestRpmRepo_RepositoryDemo(t *testing.T) {

	f, _ := ioutil.ReadFile("flufik-0.3.0-8.el8.x86_64.rpm")
	c := config.ServiceConfigBuilder{
		RpmRepositoryName: "test",
		RootRepoPath:      "/tmp",
	}
	rd := NewRpmBuilder(&c)
	err := rd.Repository(f, "flufik-0.3.0-8.el8.x86_64.rpm")
	if err != nil {
		fmt.Println(err)
	}
}
