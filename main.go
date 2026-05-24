package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	app := fiber.New()

	app.Get(
		"/*",
		static.New("./public"),
	)

	var logWriter io.Writer = os.Stdout

	app.Post(
		"/submit-consult",
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
			_, errWrite := io.WriteString(logWriter, payload)
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

	log.Fatal(app.Listen(":3000"))
}
