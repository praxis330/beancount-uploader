package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var (
	InvalidSenderError  = errors.New("You have supplied an invalid sender.")
	NoSenderError       = errors.New("Please supply a sender.")
	InternalServerError = errors.New("An internal server error has occurred.")
)

func SenderAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var RequestBody Request

		err := c.BindJSON(&RequestBody)

		if err == nil {

			c.Set("RequestBody", RequestBody)

			if RequestBody.Sender != os.Getenv("SENDER") {
				c.AbortWithError(http.StatusNotAcceptable, InvalidSenderError).SetType(gin.ErrorTypePublic)
			}
		} else {
			c.AbortWithError(http.StatusBadRequest, NoSenderError).SetType(gin.ErrorTypePublic)
		}

		c.Next()
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastError := c.Errors.Last()
		statusCode := c.Writer.Status()

		if lastError != nil {
			if lastError.IsType(gin.ErrorTypePublic) {
				c.JSON(statusCode, gin.H{
					"Error": lastError.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": InternalServerError.Error(),
			})
			return
		}

	}
}
