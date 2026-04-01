# copia-cli issue create

## copia-cli issue create

Create an issue

```
copia-cli issue create [flags]
```

### Examples

```
  copia issue create --title "Fix sensor mapping" --label bug
  copia issue create --title "Add feature" --body "Description here"
```

### Options

```
  -b, --body string     Issue body
  -h, --help            help for create
  -l, --label strings   Add labels
  -t, --title string    Issue title (required)
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli issue](copia-cli_issue.md)	 - Manage issues

