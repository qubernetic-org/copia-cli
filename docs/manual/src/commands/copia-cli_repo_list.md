# copia-cli repo list

## copia-cli repo list

List repositories

```
copia-cli repo list [flags]
```

### Examples

```
  copia repo list
  copia repo list --org my-org
  copia repo list --json fullName,description
```

### Options

```
  -h, --help           help for list
      --json strings   Output JSON with selected fields: [full_name description private updated_at]
  -L, --limit int      Maximum number of repositories (default 30)
  -o, --org string     List repositories for an organization
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli repo](copia-cli_repo.md)	 - Manage repositories

