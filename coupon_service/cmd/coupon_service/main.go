package main

import (
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfg  = config.New()
	repo = memdb.New()
)

func main() {
	serviceInstance := service.New(repo)
	apiInstance := api.New(cfg.API, serviceInstance)

	go func() {
		apiInstance.Start()
	}()

	fmt.Printf("🚀 Starting Coupon service server on %v. Press Ctrl + C to stop the server.\n", cfg)

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	// Notify sigChan when receiving SIGINT or SIGTERM
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan

	fmt.Println("🛑 Shutting down the server...")

	// Call the Close method of the apiInstance
	apiInstance.Close()

	fmt.Println("👋 Coupon service server stopped")
}
