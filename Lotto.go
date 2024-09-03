package main

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Lotto struct {
	Lid    int    `json:"Lid"`
	Number string `json:"Number"`
	Period int    `json:"Period"`
	Price  int    `json:"Price"`
}

func GetLotto(c *fiber.Ctx) error {
	query := `SELECT Lid, Number, Period, Price FROM Lotto`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	var Lottos []Lotto
	for rows.Next() {
		var p Lotto
		err := rows.Scan(&p.Lid, &p.Number, &p.Period, &p.Price)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		Lottos = append(Lottos, p)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// Send JSON response
	return c.JSON(Lottos)

}

func PostLotto(c *fiber.Ctx) error {
	// สร้าง seed เพียงครั้งเดียวก่อนวนลูป
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		randomNumber := rand.Intn(900000) + 100000
		query := `INSERT INTO Lotto(Number,Period, Price) VALUES (?,?,?)`
		_, err := db.Exec(query, randomNumber, 1, 80)
		if err != nil {
			return err
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteLotto(c *fiber.Ctx) error {

	query := `DELETE FROM Lotto`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
