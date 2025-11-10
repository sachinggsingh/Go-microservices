package db

import (
	"context"
	"log"
	"time"

	"github.com/sachinggsingh/e-comm/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client         *mongo.Client
	Database       *mongo.Database
	CartCollection *mongo.Collection
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) ConnectToDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	env := config.SetEnv()
	uri := env.MONGO_URL
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	d.Client = client
	d.Database = client.Database("micro-ecomm")

	d.CartCollection = d.Database.Collection("cart")

	log.Println("Connected to MongoDB!")
	return nil
}

func (d *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}
