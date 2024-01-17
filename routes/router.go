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
			stores.PATCH("/update/:storeUID", func(c *gin.Context) {
				controllers.UpdateStore(c, collection, store)
			})
			stores.DELETE("/remove", func(c *gin.Context) {
				controllers.RemoveStore(c, collection, store)
			})
		}

		items := v1.Group("/item")
		{
			items.POST("/add/:storeUID", func(c *gin.Context) {
				controllers.AddItem(c, collection, store)
			})
			items.PATCH("/update/:itemID", func(c *gin.Context) {
				controllers.UpdateItem(c, collection, store)
			})
			items.DELETE("/remove/:itemID", func(c *gin.Context) {
				controllers.RemoveItem(c, collection, store)
			})
		}
	}

	return r
}
