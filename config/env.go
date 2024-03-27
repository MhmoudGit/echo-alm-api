package config

import (
	"fmt"
	"os"

	"github.com/MhmoudGit/echo-alm-api/auth"
)

type Env struct {
	Address  string
	Postgres string
	Mongo    string
	Secret   string
	Sender   auth.EmailSender
}

// enviroment variables setup
func SetEnv() Env {
	port := os.Getenv("PORT")
	address := fmt.Sprintf("localhost:%s", port)
	mongo := os.Getenv("MONGO")
	postgresql := os.Getenv("POSTGRESQL")
	secret := os.Getenv("SECRET")
	senderEmail := os.Getenv("EMAIL")
	senderPassword := os.Getenv("EMAIL_PASSWORD")

	return Env{
		Address:  address,
		Postgres: postgresql,
		Mongo:    mongo,
		Secret:   secret,
		Sender: auth.EmailSender{
			Email:    senderEmail,
			Password: senderPassword,
		},
	}
}
