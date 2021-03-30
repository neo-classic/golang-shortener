package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/neo-classic/golang-shortener/api"
	"github.com/neo-classic/golang-shortener/repository/mongodb"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/neo-classic/golang-shortener/repository/redis"
	"github.com/neo-classic/golang-shortener/shortener"
)

var (
	dbProvider   string
	appPort      int
	redisUrl     string
	mongoUrl     string
	mongoDbName  string
	mongoTimeout int
)

const (
	redisDBProvider = "redis"
	mongoDBProvider = "mongo"
)

func main() {
	flag.IntVar(&appPort, "port", 8000, "Http Port")
	flag.StringVar(&dbProvider, "db-provider", "redis", "Default Database")
	flag.StringVar(&redisUrl, "redis-url", "", "Redis URL to connect")
	flag.StringVar(&mongoUrl, "mongo-url", "", "MongoDB URL to connect")
	flag.StringVar(&mongoDbName, "mongo-db-name", "", "MongoDB Database name")
	flag.IntVar(&mongoTimeout, "mongo-timeout", 5, "MongoDB Timeout")
	flag.Parse()

	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	return fmt.Sprintf(":%d", appPort)
}

func chooseRepo() shortener.Repository {
	switch dbProvider {
	case redisDBProvider:
		repo, err := redis.NewRedisRepository(redisUrl)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case mongoDBProvider:
		repo, err := mongodb.NewMongoRepository(mongoUrl, mongoDbName, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
