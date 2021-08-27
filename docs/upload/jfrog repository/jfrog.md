# Push packages to jfrog repository
___

# Upload RPM package
```shell
flufik push -d rhel8 -b flufik-0.0.1.rhel8.x86_64.rpm -p <JFRO_PWD> -m <SOURCE_DIR> -w jfrog -l <JFROG_REPO> -u <JFROG_USER>
```

# Upload Deb package
```shell
flufik push -a amd64 -d dev -b flufik_0.0.1_amd64.deb -p <JFRO_PWD> -m <SOURCE_DIR> -w jfrog -l <JFROG_REPO> -u <JFROG_USER>
```
><b>Info: </b> Some values if will not require changes like -c flag will take default values.
> During push to repository package, same time sha1, sha256 and md5 hashes provided.

```shell
flufik push -h
pushes any rpm to repository

Usage:
  flufik push [flags]

Flags:
  -a, --arch string        architecture example: for deb amd64, for rpm x86_64
  -c, --component string   only requires for deb packages to push (default "main")
  -d, --dist string        only required for deb packages to push
  -h, --help               help for push
  -b, --package string     package name for push
  -p, --password string    repository password
  -m, --path string        path from where take package
  -w, --provider string    jfrog|nexus|generic
  -l, --url string         repository url
  -u, --user string        repository user (must have permission to upload packages) (default ".")
```