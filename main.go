package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
)

func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := "TestSecret"

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim := token.Claims.(jwt.MapClaims)
	fmt.Println(claim)

	return c.Next()
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect database")
	}

	app := fiber.New()
	app.Use("/book", authRequired)
	// Book API
	app.Get("/book", func(c *fiber.Ctx) error {
		return getBooks(db, c)
	})
	app.Get("/book/:id", func(c *fiber.Ctx) error {
		return getBook(db, c)
	})
	app.Post("/book", func(c *fiber.Ctx) error {
		return createBook(db, c)
	})
	app.Put("/book/:id", func(c *fiber.Ctx) error {
		return updateBook(db, c)
	})
	app.Delete("/book/:id", func(c *fiber.Ctx) error {
		return deleteBook(db, c)
	})

	// User API
	app.Post("/register", func(c *fiber.Ctx) error {
		return createUser(db, c)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return loginUser(db, c)
	})

	db.AutoMigrate(&Book{}, &User{})
	fmt.Println("Migration successful")

	app.Listen(":8080")
}
