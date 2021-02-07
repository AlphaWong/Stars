{{define "layout"}}![test](https://github.com/AlphaWong/Stars/workflows/test/badge.svg)[![codecov](https://codecov.io/gh/AlphaWong/Stars/branch/master/graph/badge.svg?token=xuILexY8TD)](https://codecov.io/gh/AlphaWong/Stars)
# Stars
Do you remember what you star ?

# update
change to async request instead waterflow now.

# Run 
```sh
TOKEN=<GITHUB_TOKEN> go run ./main.go && cp -f ./out.md ./README.md
```

# GITHUB_TOKEN
```
see https://github.com/settings/tokens
```

# Github doc
```
https://docs.github.com/en/free-pro-team@latest/rest/reference/activity#list-repositories-starred-by-a-user
```
# Result
Language|⭐️|Repos
---|---|---
{{ range . }}{{.Language}}|{{.Stars}}|{{.Items}}
{{end}}{{end}}