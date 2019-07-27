# psadm - A library and command-line tool for EC2 System Manager Parameter Store

It provides the import and export features for [SSM Parameter Store](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/systems-manager-paramstore.html) via the library.

## Features

`psadm` provides the API client with the following additional feature to aws-sdk-go's one:
- easy to use
- [patrickmn/go-cache](https://github.com/patrickmn/go-cache) integration with automatic refresh

## Use-case

- use with Lambda functions
- use with daemon on initializaton

## Installation

```sh
go get -u github.com/nabeken/psadm/cmd/psadm
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
