package render_handler

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

const (
	templateDir  = "web/template/"
	templateBase = templateDir + "base.html"
)

func RenderTemplate(c *gin.Context, page string, data any) {
	tpl := template.Must(template.ParseFiles(templateBase, templateDir+page))
	err := tpl.ExecuteTemplate(c.Writer, "base", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error renderization Template"})
		return
	}
}

func Index(c *gin.Context) {
	RenderTemplate(c, "index.html", nil)
}

func About(c *gin.Context) {
	RenderTemplate(c, "about.html", nil)
}

func HomeTemplate(c *gin.Context) {
	RenderTemplate(c, "home.html", nil)
}

func FaqTemplate(c *gin.Context) {
	RenderTemplate(c, "faq.html", nil)
}
