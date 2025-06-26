package main

import (
	"djson/internal/adaptors/persistance"
	"djson/internal/config"
	"net/http"
	"time"

	userhandler "djson/internal/interfaces/input/api/rest/handler"
	"djson/internal/interfaces/input/api/rest/middleware"
	"djson/internal/interfaces/input/api/rest/routes"
	user "djson/internal/usecase"
	"djson/pkg/migrate"
	"fmt"
	"log"
	"os"
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
		log.Fatalf("failed to run migrations")

	}

	//passing this database to NewUserRepo
	userRepo := persistance.NewUserRepo(database) //this would connect the database to userRepo
	sessionRepo := persistance.NewSessionRepo(database)
	userService := user.NewUserService(userRepo, sessionRepo)
	userHandler := userhandler.NewUserHandler(userService)

	router := routes.InitRoutes(&userHandler)

	// instead of 100ms, 1s is used since all requests were returning 504
	timeoutWrappedRouter := middleware.TimeoutMiddleware(time.Millisecond*100, router)

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
