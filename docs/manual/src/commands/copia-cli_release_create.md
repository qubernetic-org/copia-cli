# copia-cli release create

## copia-cli release create

Create a release

```
copia-cli release create <tag> [flags]
```

### Examples

```
  copia release create v1.0.0 --title "Release 1.0.0" --notes "Changelog here"
  copia release create v2.0.0-rc.1 --draft --prerelease
```

### Options

```
      --draft          Create as draft
  -h, --help           help for create
  -n, --notes string   Release notes
      --prerelease     Mark as pre-release
  -t, --title string   Release title
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli release](copia-cli_release.md)	 - Manage releases

