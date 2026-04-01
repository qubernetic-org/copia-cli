# copia pr create

## copia pr create

Create a pull request

```
copia pr create [flags]
```

### Examples

```
  copia pr create --title "feat: add wrapper" --base main --head feature/wrapper
  copia pr create --title "fix: timeout" --base develop --head fix/timeout --body "Fixes #12"
```

### Options

```
      --base string    Base branch (default "main")
  -b, --body string    PR body
  -H, --head string    Head branch
  -h, --help           help for create
  -t, --title string   PR title (required)
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia pr](copia_pr.md)	 - Manage pull requests

