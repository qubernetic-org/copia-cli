# copia api

## copia api

Make an API request

### Synopsis

Make an authenticated request to the Copia/Gitea REST API.

```
copia api <path> [flags]
```

### Examples

```
  # Get authenticated user
  copia api /user

  # Create an issue
  copia api -X POST /repos/my-org/my-repo/issues --field title="Bug report"

  # Delete a repo
  copia api -X DELETE /repos/my-org/old-repo

  # Custom header
  copia api /user --header "Accept: application/json"
```

### Options

```
  -f, --field strings    Add JSON body field (key=value)
  -H, --header strings   Add HTTP header (key: value)
  -h, --help             help for api
  -X, --method string    HTTP method (default: GET, or POST if --field is used)
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia](copia.md)	 - Copia CLI — source control for industrial automation

