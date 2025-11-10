package db

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/sachinggsingh/e-comm/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client         *mongo.Client
	Database       *mongo.Database
	UserCollection *mongo.Collection
}

func NewDB() *Database {
	return &Database{}
}

func (d *Database) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	env := config.GetEnv()
	uri := env.MONGO_URL
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Not connected")
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	d.Client = client
	d.Database = client.Database("micro-ecomm")

	d.UserCollection = d.Database.Collection("user")
	// d.ProductCollection = d.Database.Collection("product")

	log.Println("Connected to MongoDB!")
	return nil
}

func (d *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}
