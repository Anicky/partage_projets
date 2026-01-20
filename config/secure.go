package config

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func SecurityMiddleware() gin.HandlerFunc {
	security := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		SSLRedirect:           false,
		ContentSecurityPolicy: "default-src 'self'",
	})

	return func(context *gin.Context) {
		err := security.Process(context.Writer, context.Request)
		if err != nil {
			context.Abort()

			return
		}

		context.Next()
	}
}
