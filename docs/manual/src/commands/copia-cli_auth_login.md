# copia-cli auth login

## copia-cli auth login

Authenticate with a Copia instance

```
copia-cli auth login [flags]
```

### Examples

```
  # Interactive login
  copia auth login

  # Non-interactive login (CI/agent)
  copia auth login --host app.copia.io --token YOUR_TOKEN
```

### Options

```
  -h, --help           help for login
      --host string    Copia instance hostname (default: app.copia.io)
      --token string   Personal access token
```

### SEE ALSO

* [copia-cli auth](copia-cli_auth.md)	 - Authenticate with Copia

