package main

import "github.com/gofiber/fiber/v2"

type Moneyadd struct {
	UserM  int `json:"UserM"`
	Wallet int `json:"Wallet"`
}

func UpdateMoney(c *fiber.Ctx) error {
	//UserM, err := strconv.Atoi(c.Params("id"))
	p := new(Moneyadd)

	if err := c.BodyParser(p); err != nil {
		return err
	}
	query := `UPDATE UserM SET Wallet= Wallet + ? WHERE UserM=?`

	_, err := db.Exec(query, p.Wallet, p.UserM)
	if err != nil {
		return err
	}

	return c.JSON(p)
}
