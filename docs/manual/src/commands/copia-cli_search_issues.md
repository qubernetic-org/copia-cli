# copia-cli search issues

## copia-cli search issues

Search issues in a repository

### Synopsis

Search issues within the current repository. Requires repo context (git remote or owner/repo argument).

```
copia-cli search issues <query> [flags]
```

### Examples

```
  copia search issues "sensor timeout"
  copia search issues bug --state closed
```

### Options

```
  -h, --help           help for issues
      --json strings   Output JSON with selected fields: [number title state]
  -L, --limit int      Maximum number of results (default 30)
  -s, --state string   Filter by state: {open|closed}
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli search](copia-cli_search.md)	 - Search across Copia

