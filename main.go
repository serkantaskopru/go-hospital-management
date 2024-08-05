package main

import (
	"context"
	"fmt"
	r "hospital-management/app"
	"hospital-management/configs"
	"hospital-management/migrations"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var redisClient *redis.Client
var ctx = context.Background()

func initRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
}
func main() {
	location, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		fmt.Println("Failed to load time zone:", err)
		return
	}

	now := time.Now().In(location)
	expiration := now.Add(45 * time.Minute)

	fmt.Println("Current time:", now)
	fmt.Println("Token expiration:", expiration)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	dbClient := configs.ConnectPostgreSQL()

	rm := migrations.PostgreDB{DB: dbClient}

	rm.MigrateHospital()
	rm.MigrateStaff()
	rm.MigrateUser()
	rm.MigrateLocations()
	rm.MigratePolyclinic()
	rm.MigrateHospitalClinic()
	rm.MigratePerson()
	rm.MigrateJobGroup()
	rm.MigrateTitle()
	rm.MigrateSession()
	rm.MigrateResetCode()
	rm.CreateDefaultData()

	initRedis()

	r.SetupRoutes(app, dbClient)
	r.SetupAuthRoutes(app, dbClient)

	e := app.Listen(":8080")
	if e != nil {
		return
	}
}
