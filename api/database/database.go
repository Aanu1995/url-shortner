package database

import (
	"os"

	"github.com/go-redis/redis/v9"
)


func CreateClient() (client *redis.Client){
	client = redis.NewClient(&redis.Options{
    Addr: os.Getenv("DB_ADDRESS"),
    Password: os.Getenv("DB_PASS"),
    DB: 0,
  })

	return
}