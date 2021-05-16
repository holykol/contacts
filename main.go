package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// App represents application with required dependencies
type App struct {
	store Store
}

// GET /contacts
func (a *App) contacts(c *fiber.Ctx) error {
	contacts, err := a.store.Contacts(c.Context())
	if err != nil {
		return fiber.NewError(500, fmt.Sprintf("failed to fetch contacts: %v", err))
	}

	return c.JSON(contacts)
}

// GET /contact/:id
func (a *App) contact(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, fmt.Sprintf("error parsing id: %v", err))
	}

	contact, err := a.store.Contact(c.Context(), uint(id))
	if err != nil {
		return fiber.NewError(404, fmt.Sprintf("failed to fetch contact: %v", err))
	}

	return c.JSON(contact)
}

// POST /contacts
func (a *App) createContact(c *fiber.Ctx) error {
	var contact Contact
	if err := c.BodyParser(&contact); err != nil {
		return fiber.NewError(400, fmt.Sprintf("error parsing request: %v", err))
	}

	if contact.Name == "" || contact.Phone == "" {
		return fiber.NewError(400, "missing name and/or phone body parameters")
	}

	if err := a.store.CreateContact(c.Context(), &contact); err != nil {
		return fiber.NewError(500, fmt.Sprintf("error creating contact: %v", err))
	}

	return c.JSON(contact)
}

// PATCH /contact/:id
func (a *App) editContact(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, fmt.Sprintf("error parsing id: %v", err))
	}

	contact, err := a.store.Contact(c.Context(), uint(id))
	if err != nil {
		return fiber.NewError(404, fmt.Sprintf("no contact with id %d", id))
	}

	if err != c.BodyParser(&contact) {
		return fiber.NewError(400, fmt.Sprintf("error parsing request body: %v", err))
	}

	if err := a.store.EditContact(c.Context(), uint(id), contact); err != nil {
		return fiber.NewError(500, fmt.Sprintf("error updating contact: %v", err))
	}

	return c.JSON(contact)
}

// DELETE /contact/:id
func (a *App) deleteContact(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, fmt.Sprintf("error parsing id: %v", err))
	}

	return a.store.DeleteContact(c.Context(), uint(id))
}

func main() {
	db, err := gorm.Open(sqlite.Open("contacts.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Contact{}); err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}

	app := App{
		store: &model{db},
	}

	srv := fiber.New()

	srv.Get("/contacts", app.contacts)
	srv.Get("/contact/:id", app.contact)
	srv.Post("/contacts", app.createContact)
	srv.Patch("/contact/:id", app.editContact)
	srv.Delete("/contact/:id", app.deleteContact)

	log.Fatal(srv.Listen(":8080"))
}
