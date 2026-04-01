# copia-cli pr view

## copia-cli pr view

View a pull request

```
copia-cli pr view <number> [flags]
```

### Examples

```
  copia pr view 7
  copia pr view 7 --json number,title,mergeable
```

### Options

```
  -h, --help           help for view
      --json strings   Output JSON with selected fields: [number title body state mergeable author base head created_at]
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli pr](copia-cli_pr.md)	 - Manage pull requests

