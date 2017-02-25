# psadm - A tool for EC2 System Manager Parameter Store

It provides the import and export features for [SSM Parameter Store](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/systems-manager-paramstore.html).

To Import from exported parameters in YAML:

```sh
psadm import [--dryrun] [--overwrite] [--default-kms-key-id=KMS-KEY-ID] exported.yml
```

To export parameters at give time in YAML to STDOUT:

```sh
psadm export [--with-decryption] [--key-prefix=PREFIX]
```

To get a parameter at give time:
```
psadm get [--at=TIME] [--with-decryption] KEY
```

## Tutorial

TBD
