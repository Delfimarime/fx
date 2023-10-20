package gin

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func TryAndCatchHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		TryAndCatch(c, func() error {
			c.Next()
			if len(c.Errors) > 0 {
				c.Abort()
				return c.Errors[0]
			}
			return nil
		})

	}
}

func ContentTypeMiddleware(allowedContentType ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.ContentType()
		if contentType == "" {
			c.Request.Header.Set("Content-Type", allowedContentType[0])
			contentType = allowedContentType[0]
		}
		if !containsType(allowedContentType, contentType) {
			SendHttpUnsupportedMediaType(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func AcceptMiddleware(allowedAcceptType ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.GetHeader("Accept")
		if accept == "" {
			c.Request.Header.Set("Accept", allowedAcceptType[0])
		}
		if accept != "" && accept != "*/*" && !containsType(allowedAcceptType, accept) {
			SendHttpNotAcceptable(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func containsType(slice []string, value string) bool {
	for _, a := range slice {
		if strings.TrimSpace(a) == strings.TrimSpace(value) {
			return true
		}
	}
	return false
}
