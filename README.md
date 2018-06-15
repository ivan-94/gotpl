# Gotpl: Template Loader and helper for 'html/template'

[![GoDoc](https://godoc.org/github.com/carney520/gotpl?status.svg)](https://godoc.org/github.com/carney520/gotpl) 
[![Build Status](https://travis-ci.org/carney520/gotpl.svg?branch=master)](https://travis-ci.org/carney520/gotpl)

## Usage

### Load Templates

```go
import (
  "github.com/gin-gonic/gin"
  "github.com/carney520/gotpl"
)

func main() {
  app := gin.Default()
  // load template from "./templates"
  tpl := gotpl.New("./templates")
  err := tpl.Load()
  exitOnError(err)
  app.SetHTMLTemplate(tpl.Template())
  ...
}
```

If the './templates' file structure like:

```shell
/posts
  index.html
  create.html
header.html
footer.html
```

It could execute template like that:

```go
  app.GET("/posts", func(ctx *gin.Context) {
    ctx.HTML(200, "posts/index.html", gin.H{})
  })

  app.GET("/posts/new", func(ctx *gin.Context) {
    ctx.HTML(200, "posts/create.html", gin.H{})
  })
```

The file path relative to template directory root will be the template name:

```shell
/posts
  index.html  # -> posts/index.html
  create.html # -> posts/create.html
header.html   # -> header.html
footer.html   # -> footer.html
```

the `/posts/index.html` content could be:

```html
{{template "header.html"}}
<div>posts/index.html</div>
{{template "footer.html"}}
```

### Helpers

* `append`: Add a variable to an array, or create an array; in the given context

  ```html
  {{append . "moreScripts" "js/jquery-ui-1.7.2.custom.min.js"}}

  {{range .moreStyles}}
      <link rel="stylesheet" type="text/css" href="/public/{{.}}">
  {{end}}
  ```
* `set`: Set a variable in the given context.

  ```html
  {{set . "title" "Basic Chat room"}}

  <h1>{{.title}}</h1>
  ```

* raw: Prints raw, unescaped, text. 

  ```html
  <div class="body">{{raw .blogBody}}</div>
  ```

## Related Project
 * [Masterminds/sprig](https://github.com/Masterminds/sprig): Useful template functions for Go templates