package main

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Lotto struct {
	Lid    int    `json:"Lid"`
	Number string `json:"Number"`
	Period int    `json:"Period"`
	Price  int    `json:"Price"`
}

type Amount struct {
	amount int `json:"amount"`
}
type GGLotto struct {
	Bid    int    `json:"Bid"`
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
	Price  int    `json:"Price"`
}

type Buylotto struct {
	Lid    int `json:"Lid"`
	UserM  int `json:"UserM"`
	Status int `json:"Status"`
}

func GetLotto(c *fiber.Ctx) error {
	query := `SELECT Lid, Number, Period, Price FROM Lotto where Status = 0`
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
		query := `INSERT INTO Lotto(Number, Period, Price,Status) VALUES (?,?,?,0)`
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
		query := `INSERT INTO Reward(Number,Reward,Price) VALUES (?,?,?)`

		if i == 0 {
			_, err := db.Exec(query, reward.Number, 1, 6000000)
			if err != nil {
				return err
			}
		} else if i == 1 {
			_, err := db.Exec(query, reward.Number, 2, 200000)
			if err != nil {
				return err
			}
		} else if i == 2 {
			_, err := db.Exec(query, reward.Number, 3, 80000)
			if err != nil {
				return err
			}
		} else if i > 2 && i <= 10 {
			_, err := db.Exec(query, reward.Number, 4, 40000)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(query, reward.Number, 5, 20000)
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
	query := `SELECT LLid, Number, Reward,Price FROM Reward `
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	var Reswards []Reward

	for rows.Next() {
		var p Reward
		err := rows.Scan(&p.LLid, &p.Number, &p.Reward, &p.Price)
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
	var x Amount
	query := `SELECT COUNT(*) as amount FROM basketlotto,Lotto WHERE basketlotto.Lid = Lotto.Lid and basketlotto.Lid = ? `
	_, err := db.Exec(query, x.amount)
	if err != nil {
		return err
	}
	if x.amount == 0 {
		query = `INSERT INTO basketlotto(Lid,UserM, Status) VALUES (?,?,?)`
		_, err = db.Exec(query, p.Lid, p.UserM, p.Status)
		if err != nil {
			return err
		}
		query = `UPDATE Lotto SET Status=1 where Lid = ?`
		_, err = db.Exec(query, p.Lid)
		if err != nil {
			return err
		}
		return c.JSON("Status : Ok")
	}

	return c.JSON("Status : have user buy or select")

}

func getBasketLotto(c *fiber.Ctx) error {
	userid, _ := strconv.Atoi(c.Params("id"))
	query := `SELECT basketlotto.Bid,basketlotto.Lid , Lotto.Number , Lotto.Period , Lotto.Price
			FROM basketlotto, Lotto, UserM
			WHERE basketlotto.Lid = Lotto.Lid
			AND basketlotto.UserM = UserM.UserM
			and basketlotto.UserM = ? 
			and basketlotto.Status = 1`
	rows, err := db.Query(query, userid)
	if err != nil {
		return err
	}
	var Lottos []GGLotto
	for rows.Next() {
		var p GGLotto
		err := rows.Scan(&p.Bid, &p.Lid, &p.Number, &p.Period, &p.Price)
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

func getSerachLotto(c *fiber.Ctx) error {
	// รับพารามิเตอร์ id ที่เป็นเลขล็อตเตอรี่
	Lottonumber := c.Params("id")
	if Lottonumber == "" {
		return c.Status(400).SendString("Invalid lottery number")
	}

	// ใช้คำสั่ง SQL สำหรับค้นหาหมายเลขล็อตเตอรี่ที่ขึ้นต้นด้วยพารามิเตอร์ id
	query := `SELECT Lid, Number, Period, Price FROM Lotto WHERE Number LIKE ? and Status = 0 `
	searchPattern := Lottonumber + "%"

	// รันคำสั่ง SQL พร้อมพารามิเตอร์ค้นหา
	rows, err := db.Query(query, searchPattern)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	// เก็บข้อมูลล็อตเตอรี่
	var Lottos []Lotto
	for rows.Next() {
		var p Lotto
		err := rows.Scan(&p.Lid, &p.Number, &p.Period, &p.Price)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		Lottos = append(Lottos, p)
	}

	// ตรวจสอบข้อผิดพลาดจากการวนลูป rows
	if err = rows.Err(); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// ส่งผลลัพธ์กลับเป็น JSON
	return c.JSON(Lottos)
}
func NumberOneReward(c *fiber.Ctx) error {
	query := `SELECT COUNT(*) as amount FROM Reward WHERE Reward = 1`
	var x Amount

	err := db.QueryRow(query).Scan(&x.amount)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return err
		}
	}

	if x.amount == 0 {
		// Query to select a random Lotto entry that is not in Reward
		query = `SELECT Lid, Number FROM Lotto WHERE Number NOT IN (SELECT Number FROM Reward) ORDER BY RAND() LIMIT 1;`
		var p LottoReward

		err = db.QueryRow(query).Scan(&p.Lid, &p.Number)
		if err != nil {
			return err
		}

		// Insert the new reward
		query = `INSERT INTO Reward(Number, Reward, Price) VALUES (?, ?, ?)`
		_, err = db.Exec(query, p.Number, 1, 6000000)
		if err != nil {
			return err
		}

		// Return success status
		return c.JSON(fiber.Map{"status": "ok",
			"Number": p.Number})
	}

	// If LLid is not 0, return a message indicating the reward exists
	return c.JSON(fiber.Map{"status": "have reward"})
}

func NumberTwoReward(c *fiber.Ctx) error {
	query := `SELECT COUNT(*) as amount FROM Reward WHERE Reward = 2`
	var x Amount

	err := db.QueryRow(query).Scan(&x.amount)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return err
		}
	}

	if x.amount == 0 {
		// Query to select a random Lotto entry that is not in Reward
		query = `SELECT Lid, Number FROM Lotto WHERE Number NOT IN (SELECT Number FROM Reward) ORDER BY RAND() LIMIT 1;`
		var p LottoReward

		err = db.QueryRow(query).Scan(&p.Lid, &p.Number)
		if err != nil {
			return err
		}

		// Insert the new reward
		query = `INSERT INTO Reward(Number, Reward, Price) VALUES (?, ?, ?)`
		_, err = db.Exec(query, p.Number, 2, 200000)
		if err != nil {
			return err
		}

		// Return success status
		return c.JSON(fiber.Map{"status": "ok",
			"Number": p.Number})
	}

	// If LLid is not 0, return a message indicating the reward exists
	return c.JSON(fiber.Map{"status": "have reward"})
}

func NumberThreeReward(c *fiber.Ctx) error {
	query := `SELECT COUNT(*) as amount FROM Reward WHERE Reward = 3`
	var x Amount

	err := db.QueryRow(query).Scan(&x.amount)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return err
		}
	}

	if x.amount == 0 {
		// Query to select a random Lotto entry that is not in Reward
		query = `SELECT Lid, Number FROM Lotto WHERE Number NOT IN (SELECT Number FROM Reward) ORDER BY RAND() LIMIT 1;`
		var p LottoReward

		err = db.QueryRow(query).Scan(&p.Lid, &p.Number)
		if err != nil {
			return err
		}

		// Insert the new reward
		query = `INSERT INTO Reward(Number, Reward, Price) VALUES (?, ?, ?)`
		_, err = db.Exec(query, p.Number, 3, 800000)
		if err != nil {
			return err
		}

		// Return success status
		return c.JSON(fiber.Map{"status": "ok",
			"Number": p.Number})
	}

	// If LLid is not 0, return a message indicating the reward exists
	return c.JSON(fiber.Map{"status": "have reward"})

}

func NumberFourReward(c *fiber.Ctx) error {

	query := `SELECT COUNT(*) as amount FROM Reward WHERE Reward = 4`
	var x Amount

	err := db.QueryRow(query).Scan(&x.amount)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return err
		}
	}

	if x.amount == 0 {
		query = `SELECT Lid, Number FROM Lotto WHERE Number NOT IN (SELECT Number FROM Reward) ORDER BY RAND() LIMIT 8;`
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

		for _, reward := range Reswards {
			query := `INSERT INTO Reward(Number,Reward,Price) VALUES (?,?,?)`
			_, err := db.Exec(query, reward.Number, 4, 40000)
			if err != nil {
				return err
			}

		}
		return c.JSON(fiber.Map{"status": "ok",
			"Number": Reswards})
	}
	return c.JSON(fiber.Map{"status": "have reward"})
}

func NumberFiveReward(c *fiber.Ctx) error {
	query := `SELECT COUNT(*) as amount FROM Reward WHERE Reward = 5`
	var x Amount

	err := db.QueryRow(query).Scan(&x.amount)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return err
		}
	}

	if x.amount == 0 {
		query = `SELECT Lid, Number FROM Lotto WHERE Number NOT IN (SELECT Number FROM Reward) ORDER BY RAND() LIMIT 8;`
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

		for _, reward := range Reswards {
			query := `INSERT INTO Reward(Number,Reward,Price) VALUES (?,?,?)`
			_, err := db.Exec(query, reward.Number, 5, 20000)
			if err != nil {
				return err
			}

		}
		return c.JSON(fiber.Map{"status": "ok",
			"Number": Reswards})
	}
	return c.JSON(fiber.Map{"status": "have reward"})
}

func randomRewardLottobuy(c *fiber.Ctx) error {
	query := `SELECT Lid , Number FROM Lotto Where Lid IN (SELECT Lid FROM basketlotto WHERE basketlotto.Status = 2) ORDER BY RAND() LIMIT 19`
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
	if len(Reswards) != 19 {
		var x = 19 - len(Reswards)
		query = `SELECT Lid , Number FROM Lotto ORDER BY RAND() LIMIT ?`
		rows, err = db.Query(query, x)
		for rows.Next() {
			var p LottoReward
			err := rows.Scan(&p.Lid, &p.Number)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}
			Reswards = append(Reswards, p)
		}
	}
	i := 0
	for _, reward := range Reswards {
		query := `INSERT INTO Reward(Number,Reward,Price) VALUES (?,?,?)`

		if i == 0 {
			_, err := db.Exec(query, reward.Number, 1, 6000000)
			if err != nil {
				return err
			}
		} else if i == 1 {
			_, err := db.Exec(query, reward.Number, 2, 200000)
			if err != nil {
				return err
			}
		} else if i == 2 {
			_, err := db.Exec(query, reward.Number, 3, 80000)
			if err != nil {
				return err
			}
		} else if i > 2 && i <= 10 {
			_, err := db.Exec(query, reward.Number, 4, 40000)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(query, reward.Number, 5, 20000)
			if err != nil {
				return err
			}
		}
		i++

	}

	// ใช้ Reswards1, Reswards2, Reswards3 ในการตอบสนอง
	return c.JSON("Status : Ok")
}
