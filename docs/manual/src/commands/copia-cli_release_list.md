# copia-cli release list

## copia-cli release list

List releases

```
copia-cli release list [flags]
```

### Examples

```
  copia release list
  copia release list --json tagName,name
```

### Options

```
  -h, --help           help for list
      --json strings   Output JSON with selected fields: [tag_name name draft prerelease published_at]
  -L, --limit int      Maximum number of releases (default 30)
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli release](copia-cli_release.md)	 - Manage releases

