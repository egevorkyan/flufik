# Flufik installation methods

# Manual installation on Linux and MacOS
```shell
#in case of tar file download
curl https://github.com/egevorkyan/flufik/releases/download/v0.1/flufik.tar.gz --output flufik.tar.gz
#in case of zip file download
curl https://github.com/egevorkyan/flufik/releases/download/v0.1/flufik.zip --output flufik.zip
#using tar 
tar xzvf flufik.tar.gz
#using unzip
unzip flufik.zip
# Move binary depending on OS version Linux or MacOS
#MacOS
sudo mv flufik-bins/darwin/flufik /usr/local/bin/
#Linux
sudo mv flufik-bins/linux/flufik /usr/local/bin/
#Remove garbage
rm -rf flufik*
```

# RPM based installation
1. Create repo file in /etc/yum.repos.d/flufik.repo
```shell
[Flufik]
name=Flufik
baseurl=https://repositoryeg.jfrog.io/artifactory/flufik-rpm
gpgkey=https://repositoryeg.jfrog.io/artifactory/flufik-rpm/repodata/repomd.xml.key
enabled=1
gpgcheck=1
```
2. Install flufik
```shell
sudo dnf install flufik -y
```
3. Remove flufik or update
```shell
#update
sudo dnf update flufik -y
#uninstall
sudo dnf remove flufik -y
```

# Debian based installation
1. Add public key to apt-key
```shell
#Add from one of the location below
wget -qO - https://raw.githubusercontent.com/egevorkyan/repopubkey/main/repo-public.asc | sudo apt-key add -
wget -qO - https://repositoryeg.jfrog.io/artifactory/example-repo-local/pgp-pub.pgp | sudo apt-key add -
```
2. Add repository to /etc/apt/sources.list
```shell
sudo sh -c "echo 'deb https://repositoryeg.jfrog.io/artifactory/flufik-debian flufik main' >> /etc/apt/sources.list"
```
3. Install flufik
```shell
sudo apt-get install flufik
```
4. Remove flufik or update
```shell
#update
sudo apt-get upgrade flufik
#remove
sudo apt-get remove flufik
```