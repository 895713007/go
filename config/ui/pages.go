package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func servePages(e *gin.Engine) {
	e.LoadHTMLGlob("./templates/*")
	g := e.Group("/", gin.BasicAuth(gin.Accounts{
		"admin":   "admin",
	}))

	g.GET("/", pageIndex)
	g.GET("/list", pageList)
	g.Any("/edit", pageEdit)
}

func pageList(c *gin.Context) {
	c.HTML(http.StatusOK, "list.html", map[string]interface{}{
		"kk":"vv",
	})
}

func pageEdit(c *gin.Context) {
	c.HTML(http.StatusOK, "edit.html", map[string]interface{}{
		"kk":"vv",
	})
}

func pageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"kk":"vv",
	})
}
