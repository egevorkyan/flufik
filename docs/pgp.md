# PGP key gen

Flufik can generate pgp key, passphrase protected

# Arguments
| Name  | Description  | Type  | Mandatory |
|-------|--------------|-------|---------------|
| Name  | PGP name |  string |N   |   |
| Email  | Email address   | string  | N  |   |
| Comment  | PGP comment  | string  | N  |   |
| keyType  | PGP key type  | string  | N  |   |
| passphrase  | PGP password   | string  | Y  |   |
| bits  | PGP bits  | string  | N  |   |

> Note: only when keyType is rsa, bits will be used, otherwise not necessary
> All aguments will be autogenerated if not provided by user, except <b>passphrase</b>

# Usage
```shell
#other field will be autogenerated
flufik pgp -p PassW@rd
#all arguments provided
flufik pgp -n test -e test@test.com -c testing pgp gen -k rsa -p PassW@rd -b 2048
```