package main

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserM    int    `json:"UserM"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Wallet   int    `json:"Wallet"`
}

type Postuserall struct {
	UserM    int    `json:"UserM"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Wallet   int    `json:"Wallet"`
	Type     int    `json:"Type"`
}
type Resuserall struct {
	UserM  int    `json:"UserM"`
	Name   string `json:"Name"`
	Email  string `json:"Email"`
	Wallet int    `json:"Wallet"`
	Type   int    `json:"Type"`
}
type UserUpdate struct {
	UserM int    `json:"UserM"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
}
type ResUser struct {
	UserM  int    `json:"UserM"`
	Name   string `json:"Name"`
	Email  string `json:"Email"`
	Wallet int    `json:"Wallet"`
}
type UserLogin struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type BuyLottos struct {
	UserM  int   `json:"UserM"`
	Lid    []int `json:"Lid"`
	Wallet int   `json:"Wallet"`
}

type UserchcekReward struct {
	Lid         int    `json:"Lid"`
	Price       int    `json:"Price"`
	Number      string `json:"Number"`
	Period      int    `json:"Period"`
	Reward      int    `json:"Reward"`
	Rewardprice int    `json:"Rewardprice"`
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
func GetUserByid(c *fiber.Ctx) error {
	userid, _ := strconv.Atoi(c.Params("id"))
	rows := db.QueryRow(`SELECT UserM, Name, Email , Wallet FROM UserM where UserM = ?`, userid)

	var p ResUser
	err := rows.Scan(&p.UserM, &p.Name, &p.Email, &p.Wallet)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(p)
}

func PostUser(c *fiber.Ctx) error {

	p := new(User)
	if err := c.BodyParser(p); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO UserM(Name, Email, Password, Wallet,Type) VALUES (?, ?, ?,?,0)`

	result, err := db.Exec(query, p.Name, p.Email, hashedPassword, p.Wallet)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.Password = string(hashedPassword)
	p.UserM = int(id)
	P := new(ResUser)
	P.Name = p.Name
	P.Email = p.Email
	P.Wallet = p.Wallet

	return c.JSON(P)
}

func UpdateUser(c *fiber.Ctx) error {
	//UserM, err := strconv.Atoi(c.Params("id"))
	p := new(UserUpdate)

	if err := c.BodyParser(p); err != nil {
		return err
	}
	query := `UPDATE UserM SET Name=?,Email=? WHERE UserM=?`

	_, err := db.Exec(query, p.Name, p.Email, p.UserM)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

func LoginUser(c *fiber.Ctx) error {
	p := new(UserLogin)
	if err := c.BodyParser(p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	// Query the database for the user by email
	row := db.QueryRow(`SELECT UserM, Name, Email, Password, Wallet , Type FROM UserM WHERE Email = ?`, p.Email)
	// Create a User instance to hold the queried data
	P := new(Postuserall)
	err := row.Scan(&P.UserM, &P.Name, &P.Email, &P.Password, &P.Wallet, &P.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return an error if no user is found
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		// Handle other errors (e.g., database issues)
		return fiber.NewError(fiber.StatusInternalServerError, "Database query error")
	}
	// Verify the provided password against the hashed password
	if CheckPassword(P.Password, p.Password) {
		U := new(Resuserall)
		U.UserM = P.UserM
		U.Name = P.Name
		U.Email = P.Email
		U.Wallet = P.Wallet
		U.Type = P.Type
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Login successful",
			"user":    U,
		})
	} else {
		return c.JSON("Invalid email or password")
	}
}

func Userbuylotto(c *fiber.Ctx) error {
	p := new(BuyLottos)
	if err := c.BodyParser(p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}
	for i := 0; i < len(p.Lid); i++ {
		query := `
			UPDATE basketlotto
			JOIN Lotto ON basketlotto.Lid = Lotto.Lid
			JOIN UserM ON basketlotto.UserM = UserM.UserM
			SET 
				basketlotto.Status = 2,
				Lotto.Status = 2
			WHERE 
				UserM.UserM = ?
			and 
				Lotto.Lid = ? `

		_, err := db.Exec(query, p.UserM, p.Lid[i])
		if err != nil {
			return err
		}
	}
	query := `UPDATE UserM SET Wallet= Wallet - ? WHERE UserM=?`

	_, err := db.Exec(query, p.Wallet, p.UserM)
	if err != nil {
		return err
	}
	return c.JSON("Ok")
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GUserLotto(c *fiber.Ctx) error {
	userid, _ := strconv.Atoi(c.Params("id"))
	query := `SELECT Lotto.Lid, Lotto.Number, Lotto.Period, Lotto.Price 
FROM basketlotto, Lotto 
WHERE basketlotto.Lid = Lotto.Lid 
  AND basketlotto.UserM = ?
  AND basketlotto.Status = 2
  AND Lotto.Number NOT IN (SELECT Number FROM Reward);
`

	rows, err := db.Query(query, userid)
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

func DeleteLottoBasket(c *fiber.Ctx) error {
	bid, _ := strconv.Atoi(c.Params("bid"))
	query := `DELETE FROM basketlotto WHERE bid = ?`
	_, err := db.Exec(query, bid)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func UserCheckLotto(c *fiber.Ctx) error {
	userid, _ := strconv.Atoi(c.Params("id"))
	query := `SELECT Lotto.Lid, Lotto.Price , Lotto.Number,Lotto.Period,Reward.Reward,Reward.Price as Rewardprice FROM basketlotto,Reward,Lotto WHERE basketlotto.Lid = Lotto.Lid and Lotto.Number = Reward.Number and basketlotto.UserM = ? and basketlotto.Status = 2`

	rows, err := db.Query(query, userid)
	if err != nil {
		return err
	}
	var Lottos []UserchcekReward
	for rows.Next() {
		var p UserchcekReward
		err := rows.Scan(&p.Lid, &p.Price, &p.Number, &p.Period, &p.Reward, &p.Rewardprice)
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
