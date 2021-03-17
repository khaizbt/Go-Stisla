package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth(route *gin.Engine) {
	route.LoadHTMLGlob("web/view/**/*")

	// router.LoadHTMLFiles("web/pages/login.html")
	route.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pages/login", gin.H{
			"Title": "Main website",
			"Year":  "2021",
		})
	})

	route.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pages/register", gin.H{
			"data": "Masuk",
		})
	})
}
