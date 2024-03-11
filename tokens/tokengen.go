package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/comp13xo/ecommerce/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")
var userData *mongo.Collection = database.UserData(database.Client, "Users")

func TokenGenerator(email string, firstName, lastName string, uid string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return token, refreshToken, nil

}

func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "invalid token"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token has already expired"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedat", Value: updated_at})
	upsert := true
	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userData.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	if err != nil {
		log.Panic(err)
		return
	}
}
