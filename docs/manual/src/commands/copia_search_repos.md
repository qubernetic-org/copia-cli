# copia search repos

## copia search repos

Search repositories

```
copia search repos <query> [flags]
```

### Examples

```
  copia search repos plc
  copia search repos "automation controller" --json fullName,description
```

### Options

```
  -h, --help           help for repos
      --json strings   Output JSON with selected fields: [fullName description htmlUrl]
  -L, --limit int      Maximum number of results (default 30)
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia search](copia_search.md)	 - Search across Copia

