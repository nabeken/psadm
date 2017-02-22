# psadm - A tool for managing EC2 System Manager Parameter Store

Export parameters at give time in YAML:
```sh
psadm export --at <epoch time>
```

Import from exported parameters in YAML:
```sh
psadm import --dryrun exported.yml
```

Get a parameter at give time:
```
psadm get --at <epoch time>
```
