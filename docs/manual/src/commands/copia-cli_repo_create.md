# copia-cli repo create

## copia-cli repo create

Create a repository

```
copia-cli repo create <name> [flags]
```

### Examples

```
  copia repo create my-repo
  copia repo create my-repo --org my-org --private
  copia repo create my-repo --description "PLC project"
```

### Options

```
  -d, --description string   Repository description
  -h, --help                 help for create
  -o, --org string           Create in organization
      --private              Make repository private
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli repo](copia-cli_repo.md)	 - Manage repositories

