package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/comp13xo/ecommerce/database"
	"github.com/comp13xo/ecommerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ProdCollection *mongo.Collection
	UserCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		ProdCollection: prodCollection,
		UserCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddToCart(ctx, app.ProdCollection, app.UserCollection, productId, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully added to cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveItem(ctx, app.ProdCollection, app.UserCollection, productId, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully removed from cart")

	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("userId")
		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		usertId, _ := primitive.ObjectIDFromHex(userId)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usertId}}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usertId}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = cursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledCart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := database.BuyFromCart(ctx, app.UserCollection, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuy(ctx, app.ProdCollection, app.UserCollection, productId, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully placed the order")

	}
}
