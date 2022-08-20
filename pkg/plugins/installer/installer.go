package installer

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type installer struct {
	PublicURL string
	RepoName  string
}

type data struct {
	Reponame string
	URL      string
}

const (
	t = `[{{.Reponame}}]
name={{.Reponame}}
baseurl={{.URL}}
enabled=1
gpgcheck=0
priority=1`
	shell = `#!/bin/bash
DEBIAN=/etc/debian_version
REDHAT=/etc/redhat-release
PKG=$1
install() {
  if test -f "$DEBIAN"; then
    echo "Removing any previous installed ..."
    tidydeb $PKG
    echo "Adding flufik public key ..."
    sudo curl -fsSL {{.KeyUrl}} -o /etc/apt/trusted.gpg.d/flufik.asc
    echo "Adding flufik repo to APT ..."
    echo "deb {{.DebRepoUrl}} stable main" | sudo tee /etc/apt/sources.list.d/flufik.list
    echo "Updating APT ..."
    DEBIAN_FRONTEND=noninteractive sudo apt update
    echo "Installing Package ..."
    DEBIAN_FRONTEND=noninteractive sudo apt install -y $PKG
  fi
  if test -f "$REDHAT"; then
    tidyrpm $PKG
    echo "Adding flufik repo to YUM ..."
    sudo sh -c "echo '{{.Repo}}' > /etc/yum.repo.d/flufik.repo"
    echo "Updating YUM ..."
    sudo yum update -y
    echo "Installing Package ..."
    sudo yum install $PKG -y
  fi
}
tidydeb() {
  PKG_OK=$(dpkg-query -W --showformat='${Status}\n' $1|grep "install ok installed")
  if [ ! "" = "$PKG_OK" ]; then
    sudo apt-get -y remove $1
  fi
}
tidyrpm() {
  PKG_OK=$(rpm -q $1)
  if [ !  "$PKG_OK" == "package $1 is not installed" ]; then
    sudo yum remove $1 -y
  fi
}
install
`
)

func (i *installer) TemplateValueGererate() *data {
	url := fmt.Sprintf("https://%s/%s", i.PublicURL, i.RepoName)
	return &data{
		Reponame: i.RepoName,
		URL:      url,
	}
}

func (i *installer) GenerateRepo() (string, error) {
	var n bytes.Buffer
	d := i.TemplateValueGererate()
	templ, err := template.New("repofile").Parse(t)
	if err != nil {
		return "", fmt.Errorf("failed to generate template: %v", err)
	}
	err = templ.Execute(&n, d)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}
	return n.String(), nil
}

func (i *installer) GenerateShell(pathShell string) error {
	var n bytes.Buffer
	Repo, err := i.GenerateRepo()
	if err != nil {
		return err
	}
	d := struct {
		Repo       string
		DebRepoUrl string
		KeyUrl     string
	}{
		Repo:       Repo,
		DebRepoUrl: fmt.Sprintf("https://%s", i.PublicURL),
		KeyUrl:     fmt.Sprintf("https://%s/public/flufik_pub.pgp", i.PublicURL),
	}
	templ, err := template.New("repofile").Parse(shell)
	if err != nil {
		return fmt.Errorf("failed to generate shell template: %v", err)
	}
	err = templ.Execute(&n, d)
	if err != nil {
		return fmt.Errorf("failed to execute shell template: %v", err)
	}
	f, err := os.Create(filepath.Join(pathShell, "get.sh"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	_, err = fmt.Fprint(f, n.String())
	if err != nil {
		return err
	}
	return nil
}

func NewInstaller(repoName string, publicUrl string) *installer {
	return &installer{
		RepoName:  repoName,
		PublicURL: publicUrl,
	}
}
