package service

import (
	"github.com/minio/minio-go"
	"log"
	"os"
)

func InitClientDO() (*minio.Client, error) {

	ACCESS_KEY := os.Getenv("ACCESS_KEY")
	SECRET_KEY := os.Getenv("SECRET_KEY")
	endpoint := "fra1.digitaloceanspaces.com"
	ssl := true

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, ACCESS_KEY, SECRET_KEY, ssl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return client, nil
}
