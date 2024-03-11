package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/comp13xo/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
	ErrCantRemoveItemCart = errors.New("can't remove item cart")
)

func AddToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.Product
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productId}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func BuyFromCart(ctx context.Context, userCollection *mongo.Collection, userId string) error {
	//fetch the cart of the user
	//find the cart total
	//create an order with the items
	//added order to user collection
	//added items in the cart to order list
	//empty up the cart

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var getCartItems models.User
	var orderItem models.Order

	orderItem.OrderID = primitive.NewObjectID()
	orderItem.OrderedAt = time.Now()
	orderItem.OrderCart = make([]models.ProductUser, 0)
	orderItem.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	res, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})

	if err != nil {
		panic(err)
	}
	var getusercart []bson.M
	if err = res.All(ctx, &getusercart); err != nil {
		panic(err)
	}
	var total_price int32
	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	orderItem.Price = int(total_price)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderItem}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	usercartempty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercartempty}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil

}

func InstantBuy(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var product_details models.ProductUser
	var order_detail models.Order
	order_detail.OrderID = primitive.NewObjectID()
	order_detail.OrderedAt = time.Now().Local()
	order_detail.OrderCart = make([]models.ProductUser, 0)
	order_detail.PaymentMethod.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productId}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}
	order_detail.Price = int(*product_details.Price)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
