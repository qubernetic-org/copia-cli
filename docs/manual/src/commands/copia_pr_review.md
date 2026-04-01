# copia pr review

## copia pr review

Submit a review on a pull request

```
copia pr review <number> [flags]
```

### Examples

```
  copia pr review 7 --approve
  copia pr review 7 --request-changes --body "Please fix the tests."
  copia pr review 7 --comment --body "Looks good overall."
```

### Options

```
      --approve           Approve the PR
  -b, --body string       Review body text
      --comment           Leave a review comment
  -h, --help              help for review
      --request-changes   Request changes
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia pr](copia_pr.md)	 - Manage pull requests

