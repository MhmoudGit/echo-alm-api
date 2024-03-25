package main

import (
	"context"
	"log"

	"github.com/MhmoudGit/echo-alm-api/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	connect(env ENV)
	disconnect()
}

func database(db Database, env ENV) {
	db.connect(env)
}

type Postgre struct {
	gorm *gorm.DB
}

// postgresql database connections
func (pg *Postgre) connect(env ENV) {
	var err error
	// Open db connection
	pg.gorm, err = gorm.Open(postgres.Open(env.Postgres), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	} else {
		log.Println("Database connected successfully...")
	}
}

func (pg *Postgre) migrate() {
	// Define a slice of model structs that you want to migrate.
	modelsToMigrate := []interface{}{
		// Add more model structs here if needed.
		auth.User{},
	}
	// // AutoMigrate will create tables if they don't exist based on the model structs.
	err := pg.gorm.AutoMigrate(modelsToMigrate...)
	if err != nil {
		log.Fatalf("Error migrating database tables: %v", err)
	}
	log.Println("Tables created/updated successfully...")
	// Check if the superadmin exists, and create if not
}

func (pg *Postgre) disconnect() {
	// Close db
	dbInstance, _ := pg.gorm.DB()
	_ = dbInstance.Close()
	log.Println("Database is closed successfully...")
}

type MongoDB struct {
	MongoDatabase *mongo.Database
	Client        *mongo.Client
}

// mongo database connections
func (mdb *MongoDB) connect(env ENV) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(env.Mongo).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	mdb.Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	mdb.MongoDatabase = mdb.Client.Database("sample_mflix")
	err = mdb.MongoDatabase.Client().Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func (mdb *MongoDB) disconnect() {
	if err := mdb.Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	log.Println("database client closed")
}
