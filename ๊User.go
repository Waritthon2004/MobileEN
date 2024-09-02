package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	UserM    int    `json:"UserM"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Wallet   int    `json:"Wallet"`
}

type UserLogin struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func GetUser(c *fiber.Ctx) error {
	rows, err := db.Query(`SELECT UserM, Name, Email, Password , Wallet FROM UserM`)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

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

	result, err := db.Exec(query, p.Name, p.Email, p.Password, p.Wallet)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.UserM = int(id)

	return c.JSON(p)
}

func UpdateUser(c *fiber.Ctx) error {
	//UserM, err := strconv.Atoi(c.Params("id"))
	p := new(User)

	if err := c.BodyParser(p); err != nil {
		return err
	}
	query := `UPDATE UserM SET Name=?,Email=?,Password=?,Wallet=? WHERE UserM=?`

	_, err := db.Exec(query, p.Name, p.Email, p.Password, p.Wallet, p.UserM)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

func LoginUser(c *fiber.Ctx) error {
	// Parse the request body into the UserLogin struct
	p := new(UserLogin)
	if err := c.BodyParser(p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	row := db.QueryRow(`SELECT UserM, Name, Email, Password, Wallet FROM UserM WHERE Email = ?`, p.Email)

	P := new(User)

	err := row.Scan(&P.UserM, &P.Name, &P.Email, &P.Password, &P.Wallet)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return an error if no user is found
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		// Handle other errors (e.g., database issues)
		return fiber.NewError(fiber.StatusInternalServerError, "Database query error")
	}
	return c.JSON(P)
}
