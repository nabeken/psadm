# psadm

[![Go](https://github.com/nabeken/psadm/actions/workflows/go.yml/badge.svg)](https://github.com/nabeken/psadm/actions/workflows/go.yml)

`psadm` is a library and a command-line tool for AWS Systems Manager Parameter Store.

The command-line application provides the import and export features for [SSM Parameter Store](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/systems-manager-paramstore.html) via the library.

## Features

`psadm` provides the API client with the following additional feature to aws-sdk-go's one:
- easy to use
- [patrickmn/go-cache](https://github.com/patrickmn/go-cache) integration with automatic refresh

## Use-case

- use with Lambda functions
- use with daemon on initializaton

## Library

`v2` version supports aws-sdk-go-v2.

```sh
go get github.com/nabeken/psadm/v2
```

If you want to use the library with `aws-sdk-go`, please use v0 version.

```sh
go get github.com/nabeken/psadm
```

## Command-line installation

```sh
go get -u github.com/nabeken/psadm/v2/cmd/psadm
```

### Usage

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
