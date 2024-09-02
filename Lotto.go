package main

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
