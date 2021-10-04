# Push packages to nexus repository
___

# Upload RPM package
```shell
flufik push -w nexus -l <NEXUS_REPO_URL> -u <NEXUS_USER> -p <NEXUS_PWD> -b <EXAMPLE.rpm> -n yum  -r <NEXUS_REPO_NAME>
```

# Upload Deb package
```shell
flufik push -w nexus -l <NEXUS_REPO_URL> -u <NEXUS_USER> -p <NEXUS_PWD> -b <EXAMPLE.deb> -n apt  -r <NEXUS_REPO_NAME>
```
><b>Note: </b> If -m argument is not provided flufik will take package from $HOME/.flufik/output. If you are uploading from
> different location you can specify -m argument

```shell
flufik push -h
pushes any rpm or deb packages to repositories like nexus3 and jfrog

Usage:
  flufik push [flags]

Flags:
  -a, --arch string          architecture example: for deb amd64, for rpm x86_64
  -c, --component string     only requires for deb packages to push (default "main")
  -d, --dist string          only required for deb packages to push
  -h, --help                 help for push
  -n, --nxcomponent string   Nexus components - apt or yum
  -b, --package string       package name for push
  -p, --password string      repository password
  -m, --path string          path from where take package (default "$HOME/.flufik/output")
  -w, --provider string      jfrog|nexus|generic
  -r, --repository string    repository name for apt or yum
  -l, --url string           repository url
  -u, --user string          repository user (must have permission to upload packages) (default ".")
```