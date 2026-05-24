package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"time"

	"os"
	"runtime"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

//go:embed public/*
var embeddedFS embed.FS

func main() {
	file, errCreateFile := os.OpenFile(
		"/var/log/tara-works.log",
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

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,

		helpers.TernaryWithValueIn(
			[]int{1},
			runtime.NumCPU(),
			nil,
			bytearena.WithCounterCoreCPU(),
		),
	)
	if errCrIngestor != nil {
		log.Fatal(
			"Failed to create ingestor:",
			errCrIngestor,
		)
	}
	if ingestor == nil {
		log.Fatal(
			"Create ingestor is nil.",
		)
	}

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := arenalog.NewLogger(
		&arenalog.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: arenalog.LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		arenalog.WithTimestampRFC3339UTC(ctx),
	)
	if errCrLogger != nil {
		log.Fatal(
			"Failed to create logger:",
			errCrLogger,
		)
	}
	if logger == nil {
		log.Fatal(
			"Create logger is nil.",
		)
	}

	app := fiber.New()

	publicFS, errSubtree := fs.Sub(embeddedFS, "public")
	if errSubtree != nil {
		logger.Fatal(
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
				logger.Printf(
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

	logger.Fatal(
		app.Listen(
			":80",
			fiber.ListenConfig{
				// EnablePrefork: true,
			},
		),
	)
}
