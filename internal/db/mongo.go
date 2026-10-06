package db

import (
	"context"
	"fmt"
	"time"

	"notes-api/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	opts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("mongo connect failed: %w", err)
	}

	// Bounds the startup check so the app can't hang
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background()) // don't leak the client
		return nil, nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	return client, client.Database(cfg.MongoDB), nil
}

func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}