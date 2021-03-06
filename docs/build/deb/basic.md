# Basic packaging
________
# Option 1: Binary executable without any related documents, libraries and files
1. Prepare config file - config.yaml
```yaml
meta:
  name: flufik
  version: 0.1
  release: focal
  arch: x86_64
  maintainer: Maitainer <maitainer@example.com>
  summary: |-
    Demo Deb package
  description: |-
    Demo package which can be deployed on RPM based Linux OS
files:
  generic:
    - destination: /bin/flufik
      source: /home/demo/flufik
      owner: root
      group: root
```
2. Run flufik to package based on config file provided
```shell
flufik build -p deb -c <PATH_TO_CONFIG>/config.yaml
```
><b>Note: </b> If -d argument is not provided flufik will save package to $HOME/.flufik/output. If you want to save
> to different location you can specify -d argument

# Option 2: Inline script without any related documents, libraries and files
1. Prepare config file - config.yaml
```yaml
meta:
  name: flufik
  version: 0.1
  release: focal
  arch: x86_64
  maintainer: Maitainer <maitainer@example.com>
  summary: |-
    Demo Deb package
  description: |-
    Demo package which can be deployed on RPM based Linux OS
files:
  generic:
    - destination: /bin/flufik
      body: |-
        #!/bin/bash
        echo "I am Flufik, welcome to RPM world!!!"
      mode: 0755
      owner: root
      group: root
      mtime: 2021-08-27 11:30:00
```
2. Run flufik to package based on config file provided
```shell
flufik build -p deb -c <PATH_TO_CONFIG>/config.yaml
```
><b>Note: </b> If -d argument is not provided flufik will save package to $HOME/.flufik/output. If you want to save
> to different location you can specify -d argument

# Option 3: Inline script and executable binary without any related documents, libraries and files
1. Prepare config file - config.yaml
```yaml
meta:
  name: grouped-flufik
  version: 0.1
  release: focal
  arch: x86_64
  maintainer: Maitainer <maitainer@example.com>
  summary: |-
    Demo Deb package
  description: |-
    Demo package which can be deployed on RPM based Linux OS
files:
  generic:
    - destination: /bin/flufik
      source: /home/demo/flufik
      owner: root
      group: root
    - destination: /bin/flufik-script
      body: |-
        #!/bin/bash
        echo "I am Flufik, welcome to RPM world!!!"
      mode: 0755
      owner: root
      group: root
      mtime: 2021-08-27 11:30:00
```
2. Run flufik to package based on config file provided
```shell
flufik build -p deb -c <PATH_TO_CONFIG>/config.yaml
```
><b>Note: </b> If -d argument is not provided flufik will save package to $HOME/.flufik/output. If you want to save
> to different location you can specify -d argument

# Option 4: Full basic Inline script and executable binary without any related documents, libraries and files
1. Prepare config file - config.yaml
```yaml
meta:
  name: basic-full-flufik
  version: 0.1
  release: focal
  arch: x86_64
  maintainer: Maitainer <maitainer@example.com>
  signature:
    private_key: /pgp/private.asc
    pass_phrase: test123
  license: Apache-2.0
  url: flufik.com
  os: ubuntu
  vendor: flufik
  summary: |-
    Demo Deb package
  description: |-
    Demo package which can be deployed on RPM based Linux OS
files:
  generic:
    - destination: /bin/flufik
      source: /home/demo/flufik
      owner: root
      group: root
    - destination: /bin/flufik-script
      body: |-
        #!/bin/bash
        echo "I am Flufik, welcome to RPM world!!!"
      mode: 0755
      owner: root
      group: root
      mtime: 2021-08-27 11:30:00
preinstall:
  - echo basic-full-flufik installation started;
postinstall:
  - echo basic-full-flufik installation finished successfully;
preuninstall:
  - echo basic-full-flufik uninstalling;
postuninstall:
  - echo basic-full-flufik uninstalled successfully;
```
2. Run flufik to package based on config file provided
```shell
flufik build -p deb -c <PATH_TO_CONFIG>/config.yaml
```
><b>Note: </b> If -d argument is not provided flufik will save package to $HOME/.flufik/output. If you want to save
> to different location you can specify -d argument

><b>Info: </b> if --source-directory or -d not specified, default destination location,
> where rpm files will be saved is current location from where command executed.