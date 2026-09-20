package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dzikru07/crews-system/internal/cache"
	"github.com/dzikru07/crews-system/internal/config"
	"github.com/dzikru07/crews-system/internal/database"
	"github.com/dzikru07/crews-system/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	db     *database.Database
	cache  *cache.Redis
	router *gin.Engine
	http   *http.Server
}

func New(cfg config.Config, db *database.Database, cache *cache.Redis) *Server {
	router := gin.New()

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
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
