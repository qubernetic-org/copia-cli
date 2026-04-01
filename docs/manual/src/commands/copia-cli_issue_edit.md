# copia-cli issue edit

## copia-cli issue edit

Edit an issue

```
copia-cli issue edit <number> [flags]
```

### Examples

```
  copia issue edit 12 --title "New title"
  copia issue edit 12 --add-label bug --add-label urgent
  copia issue edit 12 --assignee john --assignee jane
  copia issue edit 12 --milestone 1
```

### Options

```
      --add-label strings   Add labels
  -a, --assignee strings    Set assignees
  -b, --body string         Set body
  -h, --help                help for edit
  -m, --milestone int       Set milestone ID
  -t, --title string        Set title
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli issue](copia-cli_issue.md)	 - Manage issues

