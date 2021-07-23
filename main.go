package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/roman-wb/websocket-mover/internal/broker"
	"github.com/roman-wb/websocket-mover/internal/client"
	"github.com/roman-wb/websocket-mover/internal/mover"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
	shutdownTimeout = 15 * time.Second
)

var (
	addr = flag.String("addr", ":8080", "http service address")
)

var upgrader = websocket.Upgrader{}

func main() {
	flag.Parse()

	// Initialize Echo
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	// Echo middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Static("/web"))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "web",
		Index: "index.html",
	}))

	// Initialize Base
	broker := broker.NewBroker()
	mover := mover.NewWorker(e.Logger, broker)
	go broker.Run(mover)

	// Initialize Handlers
	e.GET("/ws", func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		client := client.NewClient(e.Logger, broker, conn)
		broker.Register(client)
		return nil
	})

	// Initialize Server
	srv := http.Server{
		Addr:         *addr,
		Handler:      e,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Start server
	go func() {
		e.Logger.Infof("Listen server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	e.Logger.Info("Try Server shutdown")
	if err := srv.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("Stoped")
}
