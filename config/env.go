package config

import (
	"fmt"
	"os"
)

type Env struct {
	Address  string
	Postgres string
	Mongo    string
	Secret   string
}

// enviroment variables setup
func SetEnv() Env {
	port := os.Getenv("PORT")
	address := fmt.Sprintf("localhost:%s", port)
	mongo := os.Getenv("MONGO")
	postgresql := os.Getenv("POSTGRESQL")
	secret := os.Getenv("SECRET")
	return Env{
		Address:  address,
		Postgres: postgresql,
		Mongo:    mongo,
		Secret:   secret,
	}
}
