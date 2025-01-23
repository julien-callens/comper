# Comper
![Version](https://img.shields.io/badge/version-v0.1.0-yellow)

Comper is a lightweight and reusable Go package for managing HTML templates, layouts, 
and rendering dynamic content in web applications. 
It integrates seamlessly with `net/http` or the `gin` framework.

## Features

- **Template Management**: Load and manage templates efficiently from a directory.
- **Layout Support**: Easily render pages with a consistent layout.
- **Gin Middleware**: Provides middleware to simplify rendering with `gin`.
- **Dynamic Data Injection**: Merge global and page-specific data.

## Installation

Use `go get` to install the package:

```bash
go get github.com/julien-callens/comper
```

Import it in your project:

```go
import "github.com/julien-callens/comper"
```

## Getting Started

### Basic Example with `net/http`

```go
package main

import (
	"log"
	"net/http"
	"github.com/julien-callens/comper"
)

func main() {
	tm, err := comper.NewTemplateManager("./web", nil)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tm.Render(w, "pages/index.gohtml", map[string]interface{}{
			"Title": "Welcome to Comper!",
		}, false)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Example with `gin`

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/julien-callens/comper"
)

func main() {
	r := gin.Default()
	tm, err := comper.NewTemplateManager("./web", nil)
	if err != nil {
		log.Fatal(err)
	}

	r.Use(comper.ApplyLayout(tm, "layout.gohtml"))

	r.GET("/", func(c *gin.Context) {
		c.Set("content", "pages/index.gohtml")
		c.Set("data", map[string]interface{}{
			"Title": "Welcome to Comper with Gin!",
		})
	})

	log.Fatal(r.Run(":8080"))
}
```

## Directory Structure

Make sure your templates follow this structure:

```
web/                       # The directory used in NewTemplateManager
├── templates/
│   ├── layout.gohtml      # Layout template
│   └── pages/
│       ├── index.gohtml   # Page-specific template
│       └── about.gohtml   # Another page-specific template
```

## Key Functions

- **`NewTemplateManager(webDir string, globalData map[string]interface{})`**  
  Initializes a new template manager.

- **`Render(w http.ResponseWriter, tmpl string, pageData interface{}, useLayout bool)`**  
  Renders a template, optionally using a layout.

- **`ApplyLayout(templateManager *TemplateManager, layout string) gin.HandlerFunc`**  
  Gin middleware for rendering templates with layouts.

## Requirements

- Go 1.20+
- (Optional) [Gin](https://github.com/gin-gonic/gin) for middleware support.

## License

This project is open-source and available under the [MIT License](LICENSE).

## Author

Created by Julien Callens.
