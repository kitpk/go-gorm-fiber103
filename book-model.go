package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func getBooks(db *gorm.DB, c *fiber.Ctx) error {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(books)
}

func getBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(book)
}

func createBook(db *gorm.DB, c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result := db.Create(book)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(book)
}

func updateBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	book := new(Book)
	result := db.First(&book, id)
	if result.Error != nil {
		return result.Error
	}

	if err := c.BodyParser(book); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result = db.Model(&book).Updates(&book)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(book)
}

func deleteBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	result := db.Delete(&Book{}, id)
	if result.Error != nil {
		return result.Error
	}

	return c.SendString("Delete Book Successful")
}
