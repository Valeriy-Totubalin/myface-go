package main

import (
	"log"
	"os"

	"github.com/Valeriy-Totubalin/myface-go"
	"github.com/Valeriy-Totubalin/myface-go/internal/app/factories"
	"github.com/Valeriy-Totubalin/myface-go/internal/delivery/http/handler"
	"github.com/Valeriy-Totubalin/myface-go/pkg/config_manager"
	"github.com/Valeriy-Totubalin/myface-go/pkg/token_manager"
)

// @title Todo App API
// @version 1.0
// @description API server for myface application

// @host localhost:8080
// @BasePath /

// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Получаем из конфига секретный ключ для jwt
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	config, err := config_manager.GetJWTConfig(pwd + "/internal/config/jwt.json")
	if nil != err {
		log.Fatal(err.Error())
		return
	}

	tokenManager, err := token_manager.NewManager(config.SecretKey)
	if nil != err {
		log.Fatal(err.Error())
		return
	}

	serviceFactory, err := factories.NewServiceFactory()
	if nil != err {
		log.Fatal(err.Error())
		return
	}

	handlers := new(handler.Handler)
	handlers.TokenManager = tokenManager
	handlers.ServiceFactory = serviceFactory

	srv := new(myface.Server)

	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http sever: %s", err.Error())
	}
}
