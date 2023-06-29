package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zardan4/todo-app-gin"
	"github.com/zardan4/todo-app-gin/pkg/handlers"
	"github.com/zardan4/todo-app-gin/pkg/repository"
	"github.com/zardan4/todo-app-gin/pkg/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) // формат для логера

	if err := initConfiguration(); err != nil { // init configuration
		logrus.Fatalf("error occured while loading configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil { // init env
		logrus.Fatalf("error occured while loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{ // create database connection
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error occured while connecting to db: %s", err.Error())
	}

	// db
	repo := repository.NewRepository(db)
	// db -> service
	service := service.NewService(repo)
	// service -> handlers
	handlers := handlers.NewHandler(service) // повертає gin.Engine, який імплементує http.Handler, тому можемо використовувати в аргументах до srv.Run()

	srv := new(todo.Server) // створюємо сервер

	go func() { // запускаємо в горутині. graceful shutdown
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil { // запускаємо сервер з роутами
			logrus.Fatalf("error occurred while running server: %s", err.Error())
		}
	}()

	logrus.Print("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT) // при одному з цих сигналів в канал буде записуватися значення, вони читатиметься і робота сервера зупинятиметься

	<-quit

	logrus.Print("server shutdowning")

	// ці дві ф-ції гарантують graceful shutdown, тобто те, що всі операції будуть виконані, але нові прийматися не будуть
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error while shutting down server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error while shutting down database: %s", err.Error())
	}
}

func initConfiguration() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
