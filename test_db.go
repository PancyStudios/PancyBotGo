//go:build tools

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	Price       int64   `bson:"price"`
	EffectValue float64 `bson:"effect_value"`
}

func main() {
	uri := os.Getenv("mongodbUrl")
	if uri == "" {
		fmt.Println("mongodbUrl is required")
		return
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return
	}
	defer client.Disconnect(context.Background())

	coll := client.Database("PancyBot").Collection("economy_items")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc Item
		if err := cursor.Decode(&doc); err != nil {
			fmt.Println("Decode error:", err)
		} else {
			fmt.Printf("Decoded: Price=%v, EffectValue=%v\n", doc.Price, doc.EffectValue)
		}
	}
}
