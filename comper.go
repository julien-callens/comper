package comper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateManager struct {
	templates  *template.Template
	webDir     string
	layout     string
	globalData map[string]any
}

func NewTemplateManager(webDir string, globalData map[string]any) (*TemplateManager, error) {
	templates, err := indexTemplates(filepath.Join(webDir, "templates"))

	if err != nil {
		return nil, fmt.Errorf("failed to index templates: %v", err)
	}

	return &TemplateManager{
		templates:  templates,
		webDir:     webDir,
		globalData: globalData,
	}, nil
}

func (tm *TemplateManager) AddGlobalData(key string, value any) {
	tm.globalData[key] = value
}

func (tm *TemplateManager) SetLayout(layout string) error {
	if tmpl := tm.templates.Lookup(layout); tmpl == nil {
		return fmt.Errorf("template %s not found", layout)
	}

	tm.layout = layout
	return nil
}

func (tm *TemplateManager) RenderWithLayout(w http.ResponseWriter, tmpl string, pageData any, layout string) {
	err := tm.SetLayout(layout)
	if err != nil {
		log.Println("Error setting layout:", err)
		http.Error(w, "Error setting layout", http.StatusInternalServerError)
		return
	}

	tm.Render(w, tmpl, pageData, true)

	err = tm.SetLayout("")
	if err != nil && err.Error() != "template  not found" {
		log.Println("Error unsetting layout:", err)
		http.Error(w, "Error unsetting layout", http.StatusInternalServerError)
		return
	}
}

func (tm *TemplateManager) Render(w http.ResponseWriter, tmpl string, pageData any, useLayout bool) {
	data := tm.mergeData(pageData)

	if useLayout {
		if tm.layout == "" {
			log.Println("Error rendering layout template: no layout set")
			http.Error(w, "Error rendering layout template", http.StatusInternalServerError)
			return
		}

		err := tm.templates.ExecuteTemplate(w, tm.layout, map[string]any{
			"Content": tm.renderContent(tmpl, data),
			"Data":    data,
		})
		if err != nil {
			log.Println("Error rendering layout template:", err)
			http.Error(w, "Error rendering layout template", http.StatusInternalServerError)
			return
		}
	} else {
		err := tm.templates.ExecuteTemplate(w, tmpl, data)
		if err != nil {
			log.Println("Error rendering template:", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}

func (tm *TemplateManager) renderContent(tmpl string, data any) string {
	var buf strings.Builder
	if err := tm.templates.ExecuteTemplate(&buf, tmpl, data); err != nil {
		return ""
	}
	return buf.String()
}

func (tm *TemplateManager) mergeData(pageData any) map[string]any {
	mergedData := map[string]any{}
	for k, v := range tm.globalData {
		mergedData[k] = v
	}

	if pd, ok := pageData.(map[string]any); ok {
		for k, v := range pd {
			mergedData[k] = v
		}
	}

	return mergedData
}

func indexTemplates(startPath string) (*template.Template, error) {
	var templates = template.New("")

	err := filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		tmpl, err := indexTemplate(startPath, path, info, err)
		if tmpl != nil {
			templates = template.Must(templates.AddParseTree(tmpl.Name(), tmpl.Tree))
			return nil
		}
		return err
	})

	return templates, err
}

func indexTemplate(startPath string, name string, info os.FileInfo, err error) (*template.Template, error) {
	if err != nil {
		return nil, err
	}
	if !info.IsDir() && filepath.Ext(name) == ".gohtml" {
		rel, err := filepath.Rel(startPath, name)
		if err != nil {
			return nil, err
		}

		rel = filepath.ToSlash(rel)

		content, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New(rel).Parse(string(content))
		if err != nil {
			return nil, err
		}

		return tmpl, err
	}
	return nil, err
}

// ApplyLayout middleware
func ApplyLayout(templateManager *TemplateManager, layout string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		tmpl, ok := c.Get("content")
		if !ok {
			c.String(500, "content not set")
			return
		}

		data, _ := c.Get("data")

		templateManager.RenderWithLayout(c.Writer, tmpl.(string), data, layout)
	}
}
