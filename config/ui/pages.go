package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/mytokenio/go_sdk/config/ui/pongin"
	"github.com/flosch/pongo2"
	"github.com/mytokenio/go_sdk/config/driver"
	"github.com/mytokenio/go_sdk/config"
)

type H = pongo2.Context
var d driver.Driver

func init() {
	//d = driver.NewHttpDriver(
	//	driver.Host("http://xxx.com"),
	//	driver.Timeout(time.Second * 10),
	//)
	d = driver.NewMockDriver()
}

func getHandler() *gin.Engine {
	e := gin.Default()

	//use pongo2 for jinja like template
	e.HTMLRender = pongin.New(pongin.RenderOptions{
		TemplateDir: "./templates",
		ContentType: "text/html; charset=utf-8",
	})
	pongo2.DefaultSet.Globals = pongo2.Context{
		"error": fError,
		"msg":   fMsg,
	}

	//serve static files
	e.Static("/static", "./static")

	//basic http auth
	//g := e.Group("/", gin.BasicAuth(gin.Accounts{
	//	"admin":   "admin",
	//}))
	g := e.Group("/")

	g.GET("/", pageIndex)
	g.GET("/list", pageList)
	g.GET("/add", pageAdd)
	g.GET("/edit/:key", pageEdit)
	g.POST("/edit", pagePostEdit)

	return e
}

func pageList(c *gin.Context) {
	vals, err := d.List()
	if err != nil {
		//pageError(c, err.Error())
		//return
	}

	flashMsg("demo list")
	demos := []string{"user", "message_center", "test"}

	for _, s := range demos {
		vals = append(vals, &driver.Value{
			K: config.DefaultServicePrefix + s,
			V: []byte("test"),
		})
	}
	c.HTML(http.StatusOK, "list.html", H{
		"vals": vals,
	})
}

func pageAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "edit.html", H{
		"key":     "",
		"content": "",
	})
}

func pageEdit(c *gin.Context) {
	key := c.Param("key")
	value, err := d.Get(key)
	if err != nil {
		//pageError(c, err.Error())
		//return
		value = driver.NewValue(key, []byte(""))
		flashMsg("key not exists, submit to create")
	}

	c.HTML(http.StatusOK, "edit.html", H{
		"key":   key,
		"value": value,
	})
}

func pagePostEdit(c *gin.Context) {
	key := c.DefaultPostForm("key", c.Param("key"))
	if key == "" {
		pageError(c, "key error")
		return
	}

	content := c.PostForm("content")
	if content == "" {
		pageError(c, "content empty")
		return
	}

	value := driver.NewValue(key, []byte(content))
	d.Set(value)

	flashMsg("update success")
	c.Redirect(302, "/edit/"+key)
}

func pageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", H{})
}

func pageError(c *gin.Context, error string) {
	fError.Value = error
	c.HTML(http.StatusOK, "error.html", H{})
}

