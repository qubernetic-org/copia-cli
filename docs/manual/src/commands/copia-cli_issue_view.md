# copia-cli issue view

## copia-cli issue view

View an issue

```
copia-cli issue view <number> [flags]
```

### Examples

```
  copia issue view 12
  copia issue view 12 --json number,title,state
```

### Options

```
  -h, --help           help for view
      --json strings   Output JSON with selected fields: [number title body state author labels created_at comments]
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli issue](copia-cli_issue.md)	 - Manage issues

