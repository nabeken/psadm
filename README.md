# psadm - A tool for EC2 System Manager Parameter Store

It provides import and export features for SSM Parameter Store.

To export parameters at give time in YAML:

```sh
psadm export [--with-decryption] [--key-prefix=PREFIX]
```

To Import from exported parameters in YAML:

```sh
psadm import [--dryrun] [--default-kms-key-id=KMS-KEY-ID] exported.yml
```

Get a parameter at give time:
```
psadm get [--at=TIME] [--with-decryption] KEY
```

## Tutorial

TBD
