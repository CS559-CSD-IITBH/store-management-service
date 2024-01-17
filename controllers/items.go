package controllers

import (
	"context"
	"net/http"

	"github.com/CS559-CSD-IITBH/store-management-service/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddItem(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	storeUID := c.Param("storeUID")

	var itemData struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Available   bool    `json:"available" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&itemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeID, err := primitive.ObjectIDFromHex(storeUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store ID"})
		return
	}

	itemModel := models.Item{
		StoreID:     storeID,
		Name:        itemData.Name,
		Description: itemData.Description,
		Available:   itemData.Available,
		Price:       itemData.Price,
	}

	_, err = collection.InsertOne(context.TODO(), itemModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item added successfully"})
}

func UpdateItem(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	storeUID := c.Param("storeUID")
	itemID := c.Param("itemID")

	// Fetch store information
	var store models.Store
	storeID, err := primitive.ObjectIDFromHex(storeUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store ID"})
		return
	}

	filter := bson.D{{Key: "_id", Value: storeID}}
	if err := collection.FindOne(context.TODO(), filter).Decode(&store); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
		return
	}

	var updateData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Available   bool    `json:"available"`
		Price       float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Construct update query
	update := bson.D{}
	if updateData.Name != "" {
		update = append(update, bson.E{Key: "name", Value: updateData.Name})
	}
	if updateData.Description != "" {
		update = append(update, bson.E{Key: "description", Value: updateData.Description})
	}
	update = append(update, bson.E{Key: "available", Value: updateData.Available})
	update = append(update, bson.E{Key: "price", Value: updateData.Price})

	// Add store ID to the filter
	filter = bson.D{{Key: "store_id", Value: storeID}, {Key: "_id", Value: itemID}}

	result, err := collection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item updated successfully"})
}

func RemoveItem(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	storeUID := c.Param("storeUID")
	itemID := c.Param("itemID")

	// Fetch store information
	var store models.Store
	storeID, err := primitive.ObjectIDFromHex(storeUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store ID"})
		return
	}

	filter := bson.D{{Key: "_id", Value: storeID}}
	if err := collection.FindOne(context.TODO(), filter).Decode(&store); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
		return
	}

	// Delete item
	result, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: itemID}, {Key: "store_id", Value: storeID}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item removed successfully"})
}
