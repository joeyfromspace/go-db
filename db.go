package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var db *mongo.Database

// ConnectOptions are connection options for the database connection
type ConnectOptions struct {
	ConnectTimeout time.Duration
	DatabaseURL    string
	Database       string
}

// Initialize a new connection to the database
func Initialize(o *ConnectOptions) (*mongo.Database, error) {
	if client != nil {
		return selectDatabase(client, o.Database)
	}

	databaseURL := o.DatabaseURL
	connectTimeout := o.ConnectTimeout

	if connectTimeout == 0 {
		connectTimeout = time.Duration(10 * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	if databaseURL == "" {
		return nil, fmt.Errorf("missing database url")
	}

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(o.DatabaseURL))

	if err != nil {
		return nil, err
	}

	err = testClient(c)
	if err != nil {
		return nil, err
	}

	client = c

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return selectDatabase(client, o.Database)
}

// GetDatabase returns the instantiated database or errors if Initialize() has not been called
func GetDatabase() (*mongo.Database, error) {
	if db != nil {
		return db, nil
	}

	return nil, fmt.Errorf("no database object exists. Initialize must be called")
}

// Close resets the state, clearing the persisted client and db objects. Initialize() must be called again.
func Close() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client.Disconnect(ctx)
	}
	db = nil
	client = nil
}

func selectDatabase(c *mongo.Client, d string) (*mongo.Database, error) {
	if db != nil {
		return db, nil
	}

	database := c.Database(d)

	db = database
	return db, nil
}

func testClient(c *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return c.Ping(ctx, readpref.Primary())
}
