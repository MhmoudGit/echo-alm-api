package main

import (
	"fmt"
	"os"
)

type ENV struct {
	Address  string
	Postgres string
	Mongo    string
}

// enviroment variables setup
func SetEnv() ENV {
	port := os.Getenv("PORT")
	address := fmt.Sprintf("localhost:%s", port)
	mongo := os.Getenv("MONGO")
	postgresql := os.Getenv("POSTGRESQL")
	return ENV{
		Address:  address,
		Postgres: postgresql,
		Mongo:    mongo,
	}
}
