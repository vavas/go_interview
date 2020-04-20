package router

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	r := gin.Default()
	r.GET("/getArticles", getArticlesHandler)
	return r
}
