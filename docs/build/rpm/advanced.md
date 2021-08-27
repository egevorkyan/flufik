# Advanced packaging
________

1. Prepare config file - config.yaml
```yaml
meta:
  name: basic-full-flufik
  version: 0.1
  release: rhel8
  arch: x86_64
  maintainer: Maitainer <maitainer@example.com>
  signature:
    private_key: /pgp/private.asc
    pass_phrase: test123
  license: Apache-2.0
  url: flufik.com
  os: rhel8
  vendor: flufik
  summary: |-
    Demo RPM package
  description: |-
    Demo package which can be deployed on RPM based Linux OS
directory:
  destination: /etc/flufik
  mode: 0755
  owner: root
  group: root
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
  config:
    - destination: /etc/flufik/flufik.cfg
      body: |-
        FLUFIK="This is flufik's configuration"
        FLUFIK_SIZE=1000
    - destination: /etc/flufik/dummy.cfg
      source: /home/demo/dummy.cfg
      owner: root
      group: root
preinstall:
  - echo basic-full-flufik installation started;
  - useradd -uid 100 -gid 100 flufik
postinstall:
  - echo basic-full-flufik installation finished successfully;
preuninstall:
  - echo basic-full-flufik uninstalling;
postuninstall:
  - echo basic-full-flufik uninstalled successfully;
  - userdel flufik
```