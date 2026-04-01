# copia search issues

## copia search issues

Search issues

```
copia search issues <query> [flags]
```

### Examples

```
  copia search issues "sensor timeout"
  copia search issues bug --state closed
```

### Options

```
  -h, --help           help for issues
      --json strings   Output JSON with selected fields: [number title state repository]
  -L, --limit int      Maximum number of results (default 30)
  -s, --state string   Filter by state: {open|closed}
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia search](copia_search.md)	 - Search across Copia

