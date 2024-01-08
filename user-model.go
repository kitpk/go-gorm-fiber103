package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func createUser(db *gorm.DB, c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	hashedPassword, err :=
		bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(fiber.Map{
		"message": "Register successful",
	})
}

func loginUser(db *gorm.DB, c *fiber.Ctx) error {
	// Get User from email
	input := new(User)
	if err := c.BodyParser(input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user := new(User)
	result := db.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		return result.Error
	}

	// Compare password
	if err :=
		bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(input.Password),
		); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// pass return jwt
	// Create JWT token
	jwtSecretKey := "TestSecret"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
	})
}
