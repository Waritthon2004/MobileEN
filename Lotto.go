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

type LottoReward struct {
	Lid    int    `json:"Lid"`
	Number string `json:"Number"`
}

type Reward struct {
	LLid   int    `json:"LLid"`
	Number string `json:"Number"`
	Reward int    `json:"Reward"`
}

type Buylotto struct {
	Lid    int `json:"Lid"`
	UserM  int `json:"UserM"`
	Status int `json:"Status"`
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
	// สร้างตัวสร้างตัวเลขสุ่มใหม่
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		randomNumber := r.Intn(900000) + 100000
		query := `INSERT INTO Lotto(Number, Period, Price) VALUES (?,?,?)`
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
func randomNumberLottoReward(c *fiber.Ctx) error {
	query := `SELECT Lid , Number FROM Lotto ORDER BY RAND() LIMIT 19`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	var Reswards []LottoReward

	for rows.Next() {
		var p LottoReward
		err := rows.Scan(&p.Lid, &p.Number)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		Reswards = append(Reswards, p)
	}

	i := 0
	for _, reward := range Reswards {
		query := `INSERT INTO Reward(Number,Reward) VALUES (?,?)`

		if i == 0 {
			_, err := db.Exec(query, reward.Number, 1)
			if err != nil {
				return err
			}
		} else if i == 1 {
			_, err := db.Exec(query, reward.Number, 2)
			if err != nil {
				return err
			}
		} else if i == 2 {
			_, err := db.Exec(query, reward.Number, 3)
			if err != nil {
				return err
			}
		} else if i > 2 && i <= 10 {
			_, err := db.Exec(query, reward.Number, 4)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(query, reward.Number, 5)
			if err != nil {
				return err
			}
		}
		i++

	}

	// ใช้ Reswards1, Reswards2, Reswards3 ในการตอบสนอง
	return c.JSON("Status : Ok")
}

func getLottoReward(c *fiber.Ctx) error {
	query := `SELECT LLid, Number, Reward FROM Reward `
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	var Reswards []Reward

	for rows.Next() {
		var p Reward
		err := rows.Scan(&p.LLid, &p.Number, &p.Reward)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		Reswards = append(Reswards, p)
	}
	var Reswards1 Reward
	var Reswards2 Reward
	var Reswards3 Reward
	var Reswards4 []Reward
	var Reswards5 []Reward

	i := 0
	for _, reward := range Reswards {
		if i == 0 {
			Reswards1 = reward
		} else if i == 1 {
			Reswards2 = reward
		} else if i == 2 {
			Reswards3 = reward
		} else if i > 2 && i <= 10 {
			Reswards4 = append(Reswards4, reward)
		} else {
			Reswards5 = append(Reswards5, reward)
		}
		i++
	}
	return c.JSON(fiber.Map{
		"Reswards1": Reswards1,
		"Reswards2": Reswards2,
		"Reswards3": Reswards3,
		"Reswards4": Reswards4,
		"Reswards5": Reswards5,
	})
}

func DeleteReward(c *fiber.Ctx) error {

	query := `DELETE FROM Reward`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func BuyLotto(c *fiber.Ctx) error {
	p := new(Buylotto)
	if err := c.BodyParser(p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}
	query := `INSERT INTO basketlotto(Lid,UserM, Status) VALUES (?,?,?)`
	_, err := db.Exec(query, p.Lid, p.UserM, p.Status)
	if err != nil {
		return err
	}

	return c.JSON("Status : Ok")

}
