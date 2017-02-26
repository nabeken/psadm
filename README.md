# psadm - A tool for EC2 System Manager Parameter Store

It provides the import and export features for [SSM Parameter Store](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/systems-manager-paramstore.html).

## Installation

```sh
go get -u github.com/nabeken/psadm
```

## Usage

To export parameters in YAML to STDOUT:

```sh
psadm export [--key-prefix=PREFIX] > exported.yml
```

Note: All `SecureString` parameters are decrypted.

To import from exported parameters in YAML:

```sh
psadm import [--dryrun] [--skip-exist] [--overwrite] [--default-kms-key-id=KMS-KEY-ID] exported.yml
```

To get a parameter at give time in YAML:
```
psadm get [--at=TIME] KEY
```

## Tutorial

TBD
