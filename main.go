package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"runtime"
	"time"

	"os"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/tudorhulban/hxhelpers"

	"github.com/tudorhulban/arenalog"
	arenafiber "github.com/tudorhulban/arenalog/arena-fiber"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

//go:embed public/*
var embeddedFS embed.FS

func main() {
	fileHTTPServer, errFileHTTP := os.OpenFile(
		"tara-works_logs.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if errFileHTTP != nil {
		log.Fatal(
			"Failed to open log file:",
			errFileHTTP,
		)
	}
	defer fileHTTPServer.Close()

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		fileHTTPServer,

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

	l, errCrLogger := arenalog.NewLogger(
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
	if l == nil {
		log.Fatal(
			"Create logger is nil.",
		)
	}

	fiberLogger := arenafiber.ALogger{
		L: l,
	}

	fiberlog.SetLogger(&fiberLogger)

	app := fiber.New()

	publicFS, errSubtree := fs.Sub(embeddedFS, "public")
	if errSubtree != nil {
		l.Fatal(
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

	app.Use(
		logger.New(
			logger.Config{
				LoggerFunc: func(c fiber.Ctx, data *logger.Data, cfg *logger.Config) error {
					fiberLogger.L.Info(
						fmt.Sprintf("%s %s %d %s",
							c.Method(),                // GET / POST etc
							c.OriginalURL(),           // full path with query
							data.Stop.Sub(data.Start), // latency
							c.IP(),                    // client IP
						),
					)

					return nil
				},
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

	fileConsult, errCreateFile := os.OpenFile(
		"tara-works_consult.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if errCreateFile != nil {
		l.Fatal(
			"Failed to open log file:",
			errCreateFile,
		)
	}
	defer fileConsult.Close()

	app.Post(
		"/submit-consult",
		submitLimiter,
		func(c fiber.Ctx) error {
			// 1. Extract form fields (handles application/x-www-form-urlencoded automatically)
			email := c.FormValue("email")
			objective := c.FormValue("objective")

			// 2. Format the inbound payload
			payload := hxhelpers.Sprintf(
				"[CONSULT_SUBMIT] Email: %s | Objective: %s\n",
				email,
				objective,
			)

			// 3. Write directly to the io.Writer
			_, errWrite := io.WriteString(fileConsult, payload)
			if errWrite != nil {
				l.Printf(
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

	l.Fatal(
		app.Listen(
			":80",
			fiber.ListenConfig{
				// EnablePrefork: true,
			},
		),
	)
}
