package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      *string            `json:"first_name" bson:"first_name" validate:"required,min=2,max=30"`
	LastName       *string            `json:"last_name" validate:"required,min=2,max=30"`
	Password       *string            `json:"password" validate:"required,min=2,max=30"`
	Email          *string            `json:"email" `
	Phone          *string            `json:"phone" validate:"required,min=10,max=10"`
	Token          *string            `json:"token"`
	RefreshToken   *string            `json:"refresh_token"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserId         string             `json:"user_id"`
	UserCart       []ProductUser      `json:"usercart" bson:"usercart"`
	AddressDetails []Address          `json:"address" bson:"address"`
	OrderStatus    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductID   primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName *string            `json:"product_name"`
	Price       *uint              `json:"price"`
	Rating      *uint8             `json:"rating"`
	Image       *string            `json:"image"`
}

type ProductUser struct {
	ProductID   primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       *uint              `json:"price" bson:"price"`
	Rating      *uint8             `json:"rating" bson:"rating"`
	Image       *string            `json:"image" bson:"image"`
}

type Address struct {
	AddressID primitive.ObjectID `json:"_id" bson:"_id"`
	House     *string            `json:"house_name" bson:"house_name"`
	Street    *string            `json:"street" bson:"street"`
	City      *string            `json:"city" bson:"city"`
	Pincode   *string            `json:"pincode"`
}

type Order struct {
	OrderID       primitive.ObjectID `json:"_id" bson:"_id"`
	OrderCart     []ProductUser
	OrderedAt     time.Time
	Price         int
	Discount      *int
	PaymentMethod Payment
}

type Payment struct {
	Digital bool
	COD     bool
}
