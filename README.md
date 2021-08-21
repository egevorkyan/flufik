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

Flags:
  -h, --help      help for flufik
  -v, --version   version for flufik

Use "flufik [command] --help" for more information about a command.
```
```shell
host@local$ flufik -v
flufik packager version 0.1
```

# Usage
><b>Info:</b> example folder contains some basic examples how to fill configuration
> for you application package, full complex configuration examples will come soon
```shell
#Assuming configuration file is already prepared, run below to build package
#-s . means source is in same dir from where you run flufik command
flufik build -c config-deb.yaml -p deb -s .
flufik build -c config-rpm.yaml -p rpm -s .
#or this way
flufik build -c /<PATH>/config-deb.yaml -p deb -s /<PATH>
flufik build -c /<PATH>/config-rpm.yaml -p rpm -s /<PATH>
```
> Flufik can build rpm and deb packages, more packages will come soon



To easely compile and archive binaries, art cli tool was used: https://github.com/gatblau/artisan/releases/tag/v1.0
