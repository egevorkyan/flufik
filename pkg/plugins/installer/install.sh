#!/bin/bash
RHEL_REPO=/etc/yum.repo.d/flufik.repo
DEB_REPO=/etc/apt/sources.list.d/flufik.list
source /etc/os-release
PKG=$1
install() {
  case "$ID" in
    rhel | centos | fedora)
      if test -f "$RHEL_REPO"; then
        echo "Cleaning cached repos"
        sudo dnf clean all -y
        echo "Updating repos"
        sudo dnf update -y
        echo "Updating application"
        sudo dnf upgrade "$PKG" -y
      else
        echo "Adding flufik repo to YUM ..."
        echo '{{.Repo}}' | sudo tee /etc/yum.repo.d/flufik.repo
        echo "Updating YUM ..."
        sudo dnf update -y
        echo "Installing Package ..."
        sudo dnf install "$PKG" -y
      fi
      ;;
    ubuntu | debian)
      if test -f "$DEB_REPO"; then
        echo "Updating APT ..."
        DEBIAN_FRONTEND=noninteractive sudo apt update
        echo "Installing Package ..."
        DEBIAN_FRONTEND=noninteractive sudo apt install -y $PKG
      else
        echo "Adding flufik public key ..."
        sudo curl -fsSL {{.KeyUrl}} -o /etc/apt/trusted.gpg.d/flufik.asc
        echo "Adding flufik repo to APT ..."
        echo "deb {{.DebRepoUrl}} $VERSION_CODENAME main" | sudo tee /etc/apt/sources.list.d/flufik.list
        echo "Updating APT ..."
        DEBIAN_FRONTEND=noninteractive sudo apt update
        echo "Installing Package ..."
        DEBIAN_FRONTEND=noninteractive sudo apt install -y $PKG
      fi
      ;;
    *)
      echo -n "Not implemented yet"
      ;;
  esac
}
install