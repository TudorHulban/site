package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"time"

	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/static"
)

//go:embed public/*
var embeddedFS embed.FS

func main() {
	file, errCreateFile := os.OpenFile(
		"tara-works_consult.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if errCreateFile != nil {
		log.Fatal(
			"Failed to open log file:",
			errCreateFile,
		)
	}
	defer file.Close()

	app := fiber.New()

	publicFS, errSubtree := fs.Sub(embeddedFS, "public")
	if errSubtree != nil {
		log.Fatal(
			"Failed to create sub FS:",
			errSubtree,
		)
	}

	app.Use(
		"/*",
		static.New(
			"",
			static.Config{
				FS: publicFS,
			},
		),
	)

	submitLimiter := limiter.New(
		limiter.Config{
			Max:        1,                              // Allow exactly 1 request...
			Expiration: _ResubmitSeconds * time.Second, // window
			KeyGenerator: func(c fiber.Ctx) string {
				return c.IP() // Track users by their IP address
			},
			LimitReached: func(c fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(
					fiber.Map{
						"error": fmt.Sprintf(
							"Submission locked. Please wait %d minute(s) before trying again.",
							_ResubmitSeconds,
						),
					},
				)
			},
		},
	)

	app.Post(
		"/submit-consult",
		submitLimiter,
		func(c fiber.Ctx) error {
			// 1. Extract form fields (handles application/x-www-form-urlencoded automatically)
			email := c.FormValue("email")
			objective := c.FormValue("objective")

			// 2. Format the inbound payload
			payload := fmt.Sprintf(
				"[CONSULT_SUBMIT] Email: %s | Objective: %s\n",
				email,
				objective,
			)

			// 3. Write directly to the io.Writer
			_, errWrite := io.WriteString(file, payload)
			if errWrite != nil {
				log.Printf(
					"Failed to write consultation data to writer: %v",
					errWrite,
				)

				return c.Status(fiber.StatusInternalServerError).
					JSON(
						fiber.Map{
							"error": "Failed to process architectural brief",
						},
					)
			}

			// 4. Acknowledge successful receipt
			return c.SendStatus(fiber.StatusOK)
		},
	)

	log.Fatal(
		app.Listen(
			":80",
			fiber.ListenConfig{
				EnablePrefork: true,
			},
		),
	)
}
