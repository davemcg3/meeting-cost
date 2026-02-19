// api is the main HTTP server for the meeting cost calculator API.
package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/yourorg/meeting-cost/backend/go/internal/cache"
	"github.com/yourorg/meeting-cost/backend/go/internal/config"
	"github.com/yourorg/meeting-cost/backend/go/internal/container"
	"github.com/yourorg/meeting-cost/backend/go/internal/handler"
	"github.com/yourorg/meeting-cost/backend/go/internal/logger"
	"github.com/yourorg/meeting-cost/backend/go/internal/middleware"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("validate config: %v", err)
	}

	// 1. Initialize Logger
	l, err := logger.NewZapLogger(os.Getenv("ENV"))
	if err != nil {
		log.Fatalf("initialize logger: %v", err)
	}

	// 2. Initialize Code Cache
	cacheClient := cache.NewRedisCache(cfg.Cache.Addr, cfg.Cache.Password, cfg.Cache.DB)

	// 3. Initialize Database
	db, err := config.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("initialize database: %v", err)
	}

	// 4. Initialize Dependency Injection Container
	ctn, err := container.NewContainer(ctx, cfg, db, cacheClient, l)
	if err != nil {
		log.Fatalf("initialize container: %v", err)
	}

	// Run AutoMigrate in development
	if cfg.Env == "development" {
		l.Info("running auto-migration")
		if err := config.AutoMigrate(db); err != nil {
			l.Error("auto-migration failed", "error", err)
		}
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Add logging middleware
	app.Use(logger.Middleware(l))

	// 5. Initialize Handlers
	meetingHandler := handler.NewMeetingHandler(ctn.MeetingService)
	authHandler := handler.NewAuthHandler(ctn.AuthService)
	orgHandler := handler.NewOrganizationHandler(ctn.OrgService)
	consentHandler := handler.NewConsentHandler(ctn.ConsentService)
	wsHandler := handler.NewWebsocketHandler(ctn.PubSub, ctn.Logger)

	// 6. Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Websocket routes
	app.Get("/ws/meetings/:id", websocket.New(wsHandler.HandleMeetingEvents))

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "ok"})
		})

		// Public consent routes
		apiV1.Get("/consent", consentHandler.GetConsent)
		apiV1.Post("/consent", consentHandler.UpdateConsent)

		auth := apiV1.Group("/auth")
		{
			auth.Post("/register", authHandler.Register)
			auth.Post("/login", authHandler.Login)
			auth.Post("/logout", authHandler.Logout)
			auth.Post("/refresh", authHandler.RefreshToken)
			auth.Get("/me", middleware.AuthRequired(ctn.AuthService), authHandler.Me)
		}

		// Private consent routes
		apiV1.Get("/consent/history", middleware.AuthRequired(ctn.AuthService), consentHandler.GetHistory)
		apiV1.Post("/consent/sync", middleware.AuthRequired(ctn.AuthService), consentHandler.SyncConsent)

		organizations := apiV1.Group("/organizations", middleware.AuthRequired(ctn.AuthService))
		{
			organizations.Get("/", orgHandler.ListOrganizations)
			organizations.Post("/", orgHandler.CreateOrganization)
			organizations.Get("/:id", orgHandler.GetOrganization)
			organizations.Put("/:id", orgHandler.UpdateOrganization)
			organizations.Delete("/:id", orgHandler.DeleteOrganization)
			organizations.Get("/:id/members", orgHandler.GetMembers)
			organizations.Post("/:id/members", orgHandler.AddMember)
			organizations.Delete("/:id/members/:memberId", orgHandler.RemoveMember)
			organizations.Patch("/:id/members/:memberId/wage", orgHandler.UpdateMemberWage)
		}

		meetings := apiV1.Group("/meetings", middleware.AuthRequired(ctn.AuthService))
		{
			meetings.Get("/", meetingHandler.ListMeetings)
			meetings.Post("/", meetingHandler.CreateMeeting)
			meetings.Get("/:id", meetingHandler.GetMeeting)
			meetings.Post("/:id/start", meetingHandler.StartMeeting)
			meetings.Post("/:id/stop", meetingHandler.StopMeeting)
			meetings.Patch("/:id/attendees", meetingHandler.UpdateAttendeeCount)
			meetings.Get("/:id/cost", meetingHandler.GetMeetingCost)
			meetings.Delete("/:id", meetingHandler.DeleteMeeting)
		}
	}

	port := cfg.Server.Port
	if p := os.Getenv("PORT"); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			port = i
		}
	}
	l.Info("listening", "port", port)
	if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
