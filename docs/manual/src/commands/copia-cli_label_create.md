# copia-cli label create

## copia-cli label create

Create a label

```
copia-cli label create [flags]
```

### Examples

```
  copia label create --name bug --color "#e11d48"
  copia label create --name feature --color "#0969da" --description "New feature"
```

### Options

```
  -c, --color string         Label color in hex (e.g. #e11d48)
  -d, --description string   Label description
  -h, --help                 help for create
  -n, --name string          Label name (required)
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli label](copia-cli_label.md)	 - Manage labels

