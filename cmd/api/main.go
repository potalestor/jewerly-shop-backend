package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/config"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/handler"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	if err := config.Init(); err != nil {
		logrus.Fatalf("error loading config: %s\n", err.Error())
	}
}

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetString("db.postgres.port"),
		Username: viper.GetString("db.postgres.username"),
		DBName:   viper.GetString("db.postgres.dbname"),
		SSLMode:  viper.GetString("db.postgres.sslmode"),
		Password: viper.GetString("db.postgres.password"),
	})
	if err != nil {
		logrus.Fatalf("Error occurred on db initialization: %s\n", err.Error())
	}

	// Init Dependecies
	repos := repository.NewRepository(db)
	services := service.NewService(service.Dependencies{
		Repos:      repos,
		HashSalt:   viper.GetString("auth.hash_salt"),
		SigningKey: []byte(viper.GetString("auth.signing_key")),
	})
	handlers := handler.NewHandler(services)

	// Create & Run HTTP Server
	server := jewerly.NewServer()
	go func() {
		if err := server.Run(viper.GetString("port"), handlers.Init()); err != nil {
			logrus.Errorf("Error occurred while running server: %s\n", err.Error())
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logrus.Errorf("error occurred while shutting down http server: %s\n", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred while closing db connection: %s\n", err.Error())
	}
}