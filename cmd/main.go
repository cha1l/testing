package main

import (
	contester "constester-go"
	"constester-go/internal/docker"
	"constester-go/internal/handler"
	"constester-go/internal/repository"
	"constester-go/internal/service"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// docker run --rm -v $(pwd)/cpp/:/app cxx-image g++ -o program /app/main.cpp
// docker cp cpp/main.cpp test:/app/main.cpp

type Config struct {
	Addr          string `toml:"addr"`
	GinMode       string `toml:"gin_mode"`
	LogLvl        string `toml:"log_level"`
	MaxContainers int    `toml:"max_containers"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	var config Config
	if err := loadConfig(&config); err != nil {
		log.Fatal("failed to load config")
	}
	gin.SetMode(config.GinMode)

	lvl, err := log.ParseLevel(config.LogLvl)
	if err != nil {
		log.Fatal("invalid log level")
	}
	log.SetLevel(lvl)

	db, err := repository.ConnectToDB()
	if err != nil {
		log.Fatal("failed connecting to mongo db ", err)
	}
	client, err := docker.NewClientDocker()
	if err != nil {
		log.Fatal("failed connecting to docker ", err)
	}

	repo := repository.NewRepository(db)
	dock := docker.NewImages(client)
	serv := service.NewService(dock, repo)
	hand := handler.NewHandler(serv)

	server := new(contester.Server)

	if err := server.StartServer(config.Addr, hand.InitEndpoints()); err != nil {
		log.Fatal("error while handling requests ", err)
	}
}

func loadConfig(config *Config) error {
	_, err := toml.DecodeFile("configs/config.toml", config)
	if err != nil {
		return err
	}
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
