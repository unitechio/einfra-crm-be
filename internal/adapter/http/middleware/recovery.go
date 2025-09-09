
package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware recovers from any panics and writes a 500 if there was one.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic.
				log.Printf("panic: %v\n%s", r, debug.Stack())

				// Return a 500 internal server error.
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "An internal server error occurred",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}
