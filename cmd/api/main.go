package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzikru07/crews-system/internal/cache"
	"github.com/dzikru07/crews-system/internal/config"
	"github.com/dzikru07/crews-system/internal/database"
	"github.com/dzikru07/crews-system/internal/server"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		fmt.Println("configuration error:", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg.Database.URL)
	if err != nil {
		fmt.Println("database connection error:", err)
		os.Exit(1)
	}
	defer db.Close()

	redis, err := cache.Open(cfg.Redis.URL)
	if err != nil {
		fmt.Println("redis connection error:", err)
		os.Exit(1)
	}
	defer redis.Close()

	app := server.New(cfg, db, redis)

	go func() {
		if err := app.Start(); err != nil && err != context.Canceled {
			fmt.Println("server error:", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		fmt.Println("shutdown error:", err)
		os.Exit(1)
	}
}
