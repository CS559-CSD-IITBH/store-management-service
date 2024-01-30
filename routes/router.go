package routes

import (
	"os"

	"github.com/CS559-CSD-IITBH/store-management-service/controllers"
	"github.com/CS559-CSD-IITBH/store-management-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(collection *mongo.Collection, store *sessions.CookieStore) *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Content-Type", "*"}
	r.Use(cors.New(config))

	auth := middlewares.Validate(store)
	r.Use(auth)

	v1 := r.Group("/api/v1")
	{
		stores := v1.Group("/store")
		{
			stores.POST("/add", func(c *gin.Context) {
				controllers.AddStore(c, collection)
			})
			stores.PATCH("/update", func(c *gin.Context) {
				controllers.UpdateStore(c, collection)
			})
			stores.DELETE("/remove", func(c *gin.Context) {
				controllers.RemoveStore(c, collection)
			})
			stores.GET("/view", func(c *gin.Context) {
				controllers.ViewStore(c, collection)
			})
		}

		items := v1.Group("/item")
		{
			items.POST("/add", func(c *gin.Context) {
				controllers.AddItem(c, collection)
			})
			items.PATCH("/update/:itemID", func(c *gin.Context) {
				controllers.UpdateItem(c, collection)
			})
			items.DELETE("/remove/:itemID", func(c *gin.Context) {
				controllers.RemoveItem(c, collection)
			})
		}
	}

	return r
}
