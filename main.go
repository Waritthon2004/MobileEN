package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	UserM    int    `json:"user_id"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Wallet   int    `json:"Wallet"`
}

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

	// Define a route
	app.Get("/user", GetUser)
	app.Post("/user", PostUser)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

func GetUser(c *fiber.Ctx) error {
	rows, err := db.Query(`SELECT UserM, Name, Email, Password FROM UserM`)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var p User
		err := rows.Scan(&p.UserM, &p.Name, &p.Email, &p.Password, &p.Wallet)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		users = append(users, p)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// Send JSON response
	return c.JSON(users)
}

func PostUser(c *fiber.Ctx) error {

	p := new(User)
	if err := c.BodyParser(p); err != nil {
		return err
	}

	query := `INSERT INTO UserM(Name, Email, Password, Wallet) VALUES (?, ?, ?,?)`

	_, err := db.Exec(query, p.Name, p.Email, p.Password, p.Wallet)
	if err != nil {
		return err
	}

	return c.JSON(p)
}
