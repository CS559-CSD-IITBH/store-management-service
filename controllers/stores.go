package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/CS559-CSD-IITBH/store-management-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddStore(c *gin.Context, collection *mongo.Collection) {
	var storeData struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&storeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchantID, err := primitive.ObjectIDFromHex(c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	storeModel := models.Store{
		MerchantID:  merchantID,
		Name:        storeData.Name,
		Description: storeData.Description,
		Available:   false,
	}

	result, err := collection.InsertOne(context.TODO(), storeModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assuming the inserted document has an "_id" field
	storeID := result.InsertedID.(primitive.ObjectID).Hex()

	c.JSON(http.StatusOK, gin.H{"status": "success", "store_id": storeID})
}

func UpdateStore(c *gin.Context, collection *mongo.Collection) {
	merchantID, err := primitive.ObjectIDFromHex(c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{
		{Key: "merchant_id", Value: merchantID},
	}

	var storeModel models.Store
	err = collection.FindOne(context.TODO(), filter).Decode(&storeModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("err: ", "No store found for the given merchant ID")
			c.JSON(http.StatusNotFound, gin.H{"error": "No store found for the given merchant ID"})
			return
		}
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	storeID := storeModel.StoreID

	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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

	filter = bson.D{
		{Key: "_id", Value: storeID},
		{Key: "merchant_id", Value: merchantID},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Store updated successfully"})
}

func RemoveStore(c *gin.Context, collection *mongo.Collection) {
	merchantID, err := primitive.ObjectIDFromHex(c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{
		{Key: "merchant_id", Value: merchantID},
	}

	var storeModel models.Store
	err = collection.FindOne(context.TODO(), filter).Decode(&storeModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("err: ", "No store found for the given merchant ID")
			c.JSON(http.StatusNotFound, gin.H{"error": "No store found for the given merchant ID"})
			return
		}
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	storeID := storeModel.StoreID

	filter = bson.D{
		{Key: "_id", Value: storeID},
		{Key: "merchant_id", Value: merchantID},
	}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Store removed successfully"})
}

func ViewStore(c *gin.Context, collection *mongo.Collection) {
	merchantID, err := primitive.ObjectIDFromHex(c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{
		{Key: "merchant_id", Value: merchantID},
	}

	var storeModel models.Store
	err = collection.FindOne(context.TODO(), filter).Decode(&storeModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("err: ", "No store found for the given merchant ID")
			c.JSON(http.StatusNotFound, gin.H{"error": "No store found for the given merchant ID"})
			return
		}
		log.Println("err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "store": storeModel})
}
