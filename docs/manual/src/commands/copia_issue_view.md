# copia issue view

## copia issue view

View an issue

```
copia issue view <number> [flags]
```

### Examples

```
  copia issue view 12
  copia issue view 12 --json number,title,state
```

### Options

```
  -h, --help           help for view
      --json strings   Output JSON with selected fields: [number title body state author labels createdAt comments]
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia issue](copia_issue.md)	 - Manage issues

