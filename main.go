package main

import (
	"log"
	"os"
	"strconv"

	"github.com/neo-classic/golang-shortener/repository/redis"
	"github.com/neo-classic/golang-shortener/shortener"
)

func main() {
	repo := chooseRepo()
}

func chooseRepo() shortener.RepositoryRedirecter {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redis.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongodb.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
