package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	PgDatabase    *gorm.DB
	MongoDatabase *mongo.Database
	Client        *mongo.Client
}

// postgresql database connections
func (s *Storage) GormConnect(env ENV) {
	var err error
	// Open db connection
	s.PgDatabase, err = gorm.Open(postgres.Open(env.Postgres), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	} else {
		log.Println("Database connected successfully...")
	}
}

func (s *Storage) GormAutoMigrateDb() {
	// Define a slice of model structs that you want to migrate.
	modelsToMigrate := []interface{}{
		// Add more model structs here if needed.
	}
	// // AutoMigrate will create tables if they don't exist based on the model structs.
	err := s.PgDatabase.AutoMigrate(modelsToMigrate...)
	if err != nil {
		log.Fatalf("Error migrating database tables: %v", err)
	}
	log.Println("Tables created/updated successfully...")
	// Check if the superadmin exists, and create if not
}

func (s *Storage) GormClose() {
	// Close db
	dbInstance, _ := s.PgDatabase.DB()
	_ = dbInstance.Close()
	log.Println("Database is closed successfully...")
}

// mongo database connections
func (s *Storage) ConnectMongoDatabase(env ENV) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(env.Mongo).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	s.Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	s.MongoDatabase = s.Client.Database("sample_mflix")
	err = s.MongoDatabase.Client().Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func (s *Storage) DisconnectMongoDatabase() {
	if err := s.Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	log.Println("database client closed")
}
