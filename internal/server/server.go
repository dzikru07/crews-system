package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dzikru07/crews-system/internal/application"
	"github.com/dzikru07/crews-system/internal/cache"
	"github.com/dzikru07/crews-system/internal/client"
	"github.com/dzikru07/crews-system/internal/config"
	"github.com/dzikru07/crews-system/internal/credential"
	"github.com/dzikru07/crews-system/internal/database"
	"github.com/dzikru07/crews-system/internal/environment"
	"github.com/dzikru07/crews-system/internal/identity"
	"github.com/dzikru07/crews-system/internal/middleware"
	"github.com/dzikru07/crews-system/internal/tenant"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	db     *database.Database
	cache  *cache.Redis
	router *gin.Engine
	http   *http.Server

	tenantHandler      *tenant.Handler
	applicationHandler *application.Handler
	environmentHandler *environment.Handler
	clientHandler      *client.Handler
	identityHandler    *identity.Handler
	credentialHandler  *credential.Handler
}

func New(cfg config.Config, db *database.Database, cache *cache.Redis) *Server {
	router := gin.New()

	tenantRepository := tenant.NewRepository(db)
	tenantService := tenant.NewService(tenantRepository)
	tenantHandler := tenant.NewHandler(tenantService)

	applicationRepository := application.NewRepository(db)
	applicationService := application.NewService(applicationRepository)
	applicationHandler := application.NewHandler(applicationService)

	environmentRepository := environment.NewRepository(db)

	environmentService := environment.NewService(
		environmentRepository,
	)

	environmentHandler := environment.NewHandler(
		environmentService,
	)

	clientRepository := client.NewRepository(db)

	clientService := client.NewService(
		clientRepository,
	)

	clientHandler := client.NewHandler(
		clientService,
	)

	identityRepository := identity.NewRepository(db)
	identityService := identity.NewService(identityRepository)
	identityHandler := identity.NewHandler(identityService)

	credentialRepository := credential.NewRepository(db)
	credentialService := credential.NewService(
		credentialRepository,
	)
	credentialHandler := credential.NewHandler(
		credentialService,
		identityRepository,
	)

	router.Use(
		middleware.RequestID(),
		middleware.Recovery(),
		middleware.SecurityHeaders(),
	)

	server := &Server{
		cfg:    cfg,
		db:     db,
		cache:  cache,
		router: router,
		http: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
			Handler: router,
		},
		tenantHandler:      tenantHandler,
		applicationHandler: applicationHandler,
		environmentHandler: environmentHandler,
		clientHandler:      clientHandler,
		identityHandler:    identityHandler,
		credentialHandler:  credentialHandler,
	}

	server.registerRoutes()

	return server
}

func (s *Server) registerRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	s.router.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)
		defer cancel()

		if err := s.db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "not_ready",
				"database": "unavailable",
			})

			return
		}

		if err := s.cache.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"redis":  "unavailable",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ready",
			"database": "ok",
			"redis":    "ok",
		})
	})

	api := s.router.Group("/v1")

	tenants := api.Group("/tenants")

	tenants.POST("", s.tenantHandler.Create)
	// tenants.GET("/:id", s.tenantHandler.GetByID)
	tenants.GET("/:tenant_id", s.tenantHandler.GetByID)

	tenantApplications := tenants.Group("/:tenant_id/applications")

	tenantApplications.POST(
		"",
		s.applicationHandler.Create,
	)

	tenantApplications.GET(
		"/:application_id",
		s.applicationHandler.GetByID,
	)

	applicationEnvironments := tenantApplications.Group(
		"/:application_id/environments",
	)

	applicationEnvironments.POST(
		"",
		s.environmentHandler.Create,
	)

	applicationEnvironments.GET(
		"/:environment_id",
		s.environmentHandler.GetByID,
	)

	environmentClients := applicationEnvironments.Group(
		"/:environment_id/clients",
	)

	environmentClients.POST(
		"",
		s.clientHandler.Create,
	)

	environmentClients.GET(
		"/:client_id",
		s.clientHandler.GetByID,
	)

	s.router.POST(
		"/v1/tenants/:tenant_id/identities",
		s.identityHandler.Create,
	)

	s.router.GET(
		"/v1/tenants/:tenant_id/identities/:identity_id",
		s.identityHandler.GetByID,
	)

	s.router.POST(
		"/v1/tenants/:tenant_id/identities/:identity_id/credentials",
		s.credentialHandler.Create,
	)

}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
