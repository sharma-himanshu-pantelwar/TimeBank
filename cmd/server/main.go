package main

import (
	"net/http"
	"time"
	"timebank/internal/adaptors/persistance"
	"timebank/internal/config"

	"fmt"
	"log"
	"os"
	userhandler "timebank/internal/interfaces/input/api/rest/handler"
	"timebank/internal/interfaces/input/api/rest/middleware"
	"timebank/internal/interfaces/input/api/rest/routes"
	user "timebank/internal/usecase"
	"timebank/pkg/migrate"
)

func main() {
	// Connect to database
	database, err := persistance.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected to database")

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// connecting and using migrations
	migrate := migrate.NewMigrate(
		database.GetDB(),
		cwd+"/migrations",
	)

	err = migrate.RunMigrations()
	if err != nil {
		log.Fatalf("failed to run migrations\n %v", err)

	}

	//passing this database to NewUserRepo
	userRepo := persistance.NewUserRepo(database) //this would connect the database to userRepo
	sessionRepo := persistance.NewSessionRepo(database)
	userService := user.NewUserService(userRepo, sessionRepo)
	userHandler := userhandler.NewUserHandler(userService)

	router := routes.InitRoutes(&userHandler)

	//1s is to be replaced with 100ms
	timeoutWrappedRouter := middleware.TimeoutMiddleware(time.Second*1, router)

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("port not found")
	}
	port := cfg.APP_PORT

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), timeoutWrappedRouter)
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}
}
