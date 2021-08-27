# flufik
Flufik helps to pack your application

# Installation
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
# Check Installation
```shell
host@local$ flufik 

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
|                                                        |
|                /,,,,\_____________/,,,,\               |
|               |,(  )/,,,,,,,,,,,,,\(  ),|              |
|                \__,,,,___,,,,,___,,,,__/               |
|                  /,,,/(')\,,,/(')\,,,\                 |
|                 |,,,,___ _____ ___,,,,|                |
|                 |,,,/   \\o_o//   \,,,|                |
|                 |,,|       |       |,,|                |
|                 |,,|   \__/|\__/   |,,|                |
|                  \,,\     \_/     /,,/                 |
|                   \__\___________/__/                  |
|     ________________/,,,,,,,,,,,,,\________________    |
|    / \,,,,,,,,,,,,,,,,___________,,,,,,,,,,,,,,,,/ \   |
|   (   ),,,,,,,,,,,,,,/           \,,,,,,,,,,,,,,(   )  |
|    \_/____________,,/             \,,____________\_/   |
|                  /,/               \,\                 |
|                 |,|   I am Flufik   |,|                |
|                 |,|  ready to pack  |,|                |
|                 |,|  apps for Linux |,|                |
|                 |,|                 |,|                |
|                  \,\       O       /,/                 |
|                  /,,\_____________/,,\                 |
|                 /,,,,,,,,,,,,,,,,,,,,,\                |
|                /,,,,,,,,_______,,,,,,,,\               |
|               /,,,,,,,,/       \,,,,,,,,\              |
|              /,,,,,,, /         \,,,,,,,,\             |
|             /_____,,,/           \,,,_____\            |
|            //     \,/             \,/     \\           |
|            \\_____//               \\_____//           |
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Usage:
  flufik [command]

Available Commands:
  build       builds deployment rpm or deb or both packages
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  push        pushes any rpm to repository

Flags:
  -h, --help      help for flufik
  -v, --version   version for flufik
  
```
```shell
host@local$ flufik -v
flufik packager version 0.1
```

# Detailed documentation
<b>RPM</b>
- [Basic package build](docs/build/rpm/basic.md)
- [Advanced package build](docs/build/rpm/advanced.md)
- [Full config packaging options](docs/build/rpm/available%20configuration.md)

<b>DEB</b>
- [Basic package build](docs/build/deb/basic.md)
- [Advanced package build](docs/build/deb/advanced.md)
- [Full config packaging options](docs/build/deb/available%20configuration.md)

<b>Upload</b>
>Currently upload feature is implemented for jfrog repository, nexus and other
> will be available
- [JFrog Repository](docs/upload/jfrog%20repository/jfrog.md)

> Flufik can build rpm and deb packages, more packages will come soon



To easely compile and archive binaries, art cli tool was used: https://github.com/gatblau/artisan/releases/tag/v1.0
