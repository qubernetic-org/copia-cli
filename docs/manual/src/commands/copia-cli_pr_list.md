# copia-cli pr list

## copia-cli pr list

List pull requests

```
copia-cli pr list [flags]
```

### Examples

```
  copia pr list
  copia pr list --state closed
  copia pr list --json number,title,state
```

### Options

```
  -h, --help           help for list
      --json strings   Output JSON with selected fields: [number title state author base head updated_at]
  -L, --limit int      Maximum number of pull requests (default 30)
  -s, --state string   Filter by state: {open|closed|all} (default "open")
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli pr](copia-cli_pr.md)	 - Manage pull requests

