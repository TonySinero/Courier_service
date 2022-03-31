package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"stlab.itechart-group.com/go/food_delivery/courier_service/GRPC/grpcServer"
	"stlab.itechart-group.com/go/food_delivery/courier_service/GRPCC/grpcClient"
	"stlab.itechart-group.com/go/food_delivery/courier_service/controller"
	"stlab.itechart-group.com/go/food_delivery/courier_service/dao"
	"stlab.itechart-group.com/go/food_delivery/courier_service/pkg/database"
	"stlab.itechart-group.com/go/food_delivery/courier_service/server"
	"stlab.itechart-group.com/go/food_delivery/courier_service/service"
	"syscall"
)

// @title Courier Service
// @description Courier Service for Food Delivery Application
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.Println("Start...")
	databases, err := database.NewPostgresDB(database.PostgresDB{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DATABASE"),
		SSLMode:  os.Getenv("DB_SSL_MODE")})
	if err != nil {
		log.Fatal("failed to initialize dao:", err.Error())
	}
	grpcCli := grpcClient.NewGRPCClient(os.Getenv("HOST"))
	repository := dao.NewRepository(databases)
	services := service.NewService(repository, grpcCli)
	handlers := controller.NewHandler(services)
	port := os.Getenv("API_SERVER_PORT")

	serv := new(server.Server)

	go func() {
		err := serv.Run(port, handlers.InitRoutesGin())
		if err != nil {
			log.Fatalf("Error occured while running http server: %s", err.Error())
		}
	}()
	go func() {
		grpcServer.NewGRPCServer(services)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := serv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error occured while shutting down http server: %s", err.Error())
	}

}
