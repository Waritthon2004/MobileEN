package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var db *sql.DB

func main() {
	// Database connection string
	dsn := "web66_65011212075:65011212075@csmsu@tcp(202.28.34.197:3306)/web66_65011212075"

	// Connect to the database
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")

	// Initialize Fiber app
	app := fiber.New()

	// User
	app.Get("/user", GetUser)
	app.Post("/user/login", LoginUser)
	app.Post("/user", PostUser)
	app.Put("/user", UpdateUser)

	//Lotto
	app.Get("/lotto", GetLotto)
	app.Post("/lotto", PostLotto)
	app.Delete("/lotto", DeleteLotto)
	// Start the server
	log.Fatal(app.Listen(":3000"))
}
