package config

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Database holds the MongoDB client and database references
type Database struct {
    Client   *mongo.Client
    Database *mongo.Database
}

// ConnectDB establishes a connection to MongoDB and returns a Database instance.
// It configures connection pooling and sets appropriate timeouts.
func ConnectDB(uri, dbName string) (*Database, error) {
    // Create a context with timeout for the connection attempt
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Configure client options with connection pooling
    clientOptions := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(50).                       // Maximum connections in the pool
        SetMinPoolSize(10).                       // Minimum connections to maintain
        SetMaxConnIdleTime(30 * time.Second).     // Close idle connections after 30s
        SetServerSelectionTimeout(5 * time.Second) // Timeout for server selection

    // Connect to MongoDB
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    // Verify the connection by pinging the server
    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }

    log.Println("Connected to MongoDB successfully")

    return &Database{
        Client:   client,
        Database: client.Database(dbName),
    }, nil
}

// Disconnect closes the MongoDB connection gracefully.
// Call this when shutting down the application.
func (db *Database) Disconnect() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return db.Client.Disconnect(ctx)
}

// GetCollection returns a reference to the specified collection.
// Use this to perform operations on a specific collection.
func (db *Database) GetCollection(name string) *mongo.Collection {
    return db.Database.Collection(name)
}