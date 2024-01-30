package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Middleware checks if the user is authenticated as a merchant
func Validate(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Gorilla sessions
		session, err := store.Get(c.Request, "auth")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		uid, ok := session.Values["user"]
		log.Println("UID:", uid)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Pass the UID to the next handler
		c.Set("uid", uid)
		c.Next()
	}
}
