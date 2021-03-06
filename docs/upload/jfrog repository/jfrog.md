# Push packages to jfrog repository
___

# Upload RPM package
```shell
flufik push -w jfrog -d <DIST_NAME> -b <EXAMPLE.rpm> -a <x86_64|arm> -r <JFROG_REPO_NAME> -l <JFROG_REPO_URL> -u <JFROG_USER> -p <JFRO_PWD>
```

# Upload Deb package
```shell
flufik push -w jfrog -d <DIST_NAME> -b <EXAMPLE.deb> -a <amd64|arm> -r <JFROG_REPO_NAME> -l <JFROG_REPO_URL> -u <JFROG_USER> -p <JFRO_PWD>
```
><b>Info: </b> Some values if will not require changes like -c flag will take default values.
> During push to repository package, same time sha1, sha256 and md5 hashes provided.

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