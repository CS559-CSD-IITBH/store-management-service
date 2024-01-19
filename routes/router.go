package routes

import (
	"github.com/CS559-CSD-IITBH/store-management-service/controllers"
	"github.com/CS559-CSD-IITBH/store-management-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(collection *mongo.Collection, store *sessions.FilesystemStore) *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	auth := middlewares.SessionAuth(store)
	r.Use(auth)

	v1 := r.Group("/api/v1")
	{
		stores := v1.Group("/store")
		{
			stores.POST("/add", func(c *gin.Context) {
				controllers.AddStore(c, collection, store)
			})
			stores.PATCH("/update/:storeID", func(c *gin.Context) {
				controllers.UpdateStore(c, collection, store)
			})
			stores.DELETE("/remove/:storeID", func(c *gin.Context) {
				controllers.RemoveStore(c, collection, store)
			})
			stores.GET("/view/:storeID", func(c *gin.Context) {
				controllers.ViewStore(c, collection, store)
			})
		}

		items := v1.Group("/item")
		{
			items.POST("/add/:storeID", func(c *gin.Context) {
				controllers.AddItem(c, collection, store)
			})
			items.PATCH("/update/:storeID/:itemID", func(c *gin.Context) {
				controllers.UpdateItem(c, collection, store)
			})
			items.DELETE("/remove/:storeID/:itemID", func(c *gin.Context) {
				controllers.RemoveItem(c, collection, store)
			})
		}
	}

	return r
}
