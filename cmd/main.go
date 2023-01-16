package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"todolist/consts"
	"todolist/pkg/handler"
	"todolist/pkg/repository"
	"todolist/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// @title TodoApp API
// @version 1.0
// @description Todo application API

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService((repos))
	handlers := handler.NewHandler(services)

	srv := new(consts.Server)

	// горутина нужна для того, чтобы при отключении приложения закрыть обработку всех запросов к базе данных и серверу
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Todo app started...")

	// Внутри srv.Run запускает бесконечный цикл, что блокирует остановку программы, теперь этого не происходит
	// Поэтому создадим канал и добавим блокировку функции main с его помощью
	quit := make(chan os.Signal, 1)
	// Запись в канал будет происходить, когда процесс,
	// в котором выполняется приложение получить сигнал от системы типа sygterm или sigInt
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// Эта строка будет блокировать выполнение главной горутины main
	// Чтение из канала
	<-quit

	logrus.Print("Todo app shutdown")

	// Строчки ниже гарантируют нам, что все соединения будут обработаны перед закрытием приложения
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occuring on server shuting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occuring on db close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
