# copia pr merge

## copia pr merge

Merge a pull request

```
copia pr merge <number> [flags]
```

### Examples

```
  copia pr merge 7
  copia pr merge 7 --squash
  copia pr merge 7 --rebase --delete-branch
```

### Options

```
      --delete-branch   Delete branch after merge
  -h, --help            help for merge
      --merge           Merge commit (default)
      --rebase          Rebase and merge
      --squash          Squash and merge
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia pr](copia_pr.md)	 - Manage pull requests

