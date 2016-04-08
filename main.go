package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func uploadBeancountItem(c *gin.Context) {
	request := c.MustGet("RequestBody").(Request)
	var beancount BeancountItem

	err := GetBeancountItem(&request, &beancount)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	title := beancount.GetTitle()

	aws := InitAWSUploader(title)

	err = aws.Upload(&beancount)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "true",
	})
}

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(ErrorMiddleware())
	router.Use(SenderAuthMiddleware())

	router.POST("/beancount/upload", uploadBeancountItem)

	return router
}

func main() {
	router := InitRouter() 
	router.Run()
}
