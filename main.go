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
	dsn := "aemandko_Tinchai:Tinchai@tcp(119.59.96.110:3306)/aemandko_Tinchai"

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
	app.Get("/user/:id", GetUserByid)
	app.Post("/user/login", LoginUser)
	app.Post("/user", PostUser)
	app.Put("/user", UpdateUser)
	app.Put("/userbylotto", Userbuylotto)
	app.Get("/lottouser/:id", GUserLotto)
	app.Get("/userchecklotto/:id", UserCheckLotto)
	app.Get("/updatestatusreward/:bid", UpdateStatus)

	//Lotto
	app.Get("/lotto", GetLotto)
	app.Get("/reward", getLottoReward)
	app.Post("/lotto", PostLotto)
	app.Delete("/lotto", DeleteLotto)
	app.Delete("/reward", DeleteReward)
	app.Post("/buylottobasket", BuyLotto)
	app.Get("/basketlotto/:id", getBasketLotto)
	app.Get("/search/:id", getSerachLotto)
	app.Get("/addreward", randomNumberLottoReward)
	app.Get("/addrewarduserbuy", randomRewardLottobuy)

	app.Get("/addrewardOne", NumberOneReward)
	app.Get("/addrewardTwo", NumberTwoReward)
	app.Get("/addrewardThree", NumberThreeReward)
	app.Get("/addrewardFour", NumberFourReward)
	app.Get("/addrewardFive", NumberFiveReward)
	app.Delete("/deleteLotto/:bid", DeleteLottoBasket)
	//Money
	app.Put("/money", UpdateMoney)

	app.Delete("/Reset", ResetLL)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

func ResetLL(c *fiber.Ctx) error {
	query := `DELETE FROM basketlotto `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	query = `DELETE FROM Reward `
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	query = `DELETE FROM UserM WHERE Type = 0 `
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	query = `DELETE FROM Lotto `
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
