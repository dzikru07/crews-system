package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf(
					"PANIC: method=%s path=%s request_id=%s error=%v",
					c.Request.Method,
					c.Request.URL.Path,
					c.GetString("request_id"),
					err,
				)

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"error": "internal_server_error",
					},
				)
			}
		}()

		c.Next()
	}
}
