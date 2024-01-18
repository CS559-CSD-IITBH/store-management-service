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
	storeID, err := primitive.ObjectIDFromHex(c.Param("storeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID format"})
		return
	}

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

	itemModel := models.Item{
		ItemID:      primitive.NewObjectID(),
		Name:        itemData.Name,
		Description: itemData.Description,
		Available:   itemData.Available,
		Price:       itemData.Price,
	}

	// Update the store's Items field by appending the new item
	update := bson.M{
		"$push": bson.M{"items": itemModel},
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": storeID},
		update,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item added successfully"})
}

func UpdateItem(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	storeID, err := primitive.ObjectIDFromHex(c.Param("storeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID format"})
		return
	}

	itemID, err := primitive.ObjectIDFromHex(c.Param("itemID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID format"})
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

	// Construct update query using bson.D
	update := bson.D{}
	if updateData.Name != "" {
		update = append(update, bson.E{Key: "items.$.name", Value: updateData.Name})
	}
	if updateData.Description != "" {
		update = append(update, bson.E{Key: "items.$.description", Value: updateData.Description})
	}
	if updateData.Available != false { // Only update if available is provided
		update = append(update, bson.E{Key: "items.$.available", Value: updateData.Available})
	}
	if updateData.Price != 0 { // Only update if price is provided
		update = append(update, bson.E{Key: "items.$.price", Value: updateData.Price})
	}

	// Construct filter with store ID and item ID
	filter := bson.D{{Key: "_id", Value: storeID}, {Key: "items._id", Value: itemID}}

	// Use "$set" to update fields
	updateResult, err := collection.UpdateOne(
		context.TODO(),
		filter,
		bson.D{{Key: "$set", Value: update}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item updated successfully"})
}

func RemoveItem(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	storeID, err := primitive.ObjectIDFromHex(c.Param("storeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID format"})
		return
	}
	itemID, err := primitive.ObjectIDFromHex(c.Param("itemID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID format"})
		return
	}

	// Fetch store information
	var store models.Store
	filter := bson.D{{Key: "_id", Value: storeID}}
	if err := collection.FindOne(context.TODO(), filter).Decode(&store); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	// Delete item
	result, err := collection.UpdateOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: storeID}},
		bson.D{{Key: "$pull", Value: bson.D{{Key: "items", Value: bson.D{{Key: "_id", Value: itemID}}}}}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item removed successfully"})
}
