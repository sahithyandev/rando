package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rando-api/internal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loaded .env file")
	
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/page-view", internal.GetPageView)
	mux.HandleFunc("/live-users-count", internal.GetLiveUsersCount)

	internal.StartClientCleanup(30 * time.Second)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	fmt.Printf("Running on port %s\n", PORT)
	PORT = fmt.Sprintf(":%s", PORT)

	handler := cors.Default().Handler(mux)
	server := &http.Server {
		Addr: PORT,
		Handler: handler,
	}
	serverShutdown := make(chan struct{})
	
	go func() {
		fmt.Println("Server is running.")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error starting server: %s\n", err)
		}
		close(serverShutdown)
	}()
	
	go func() {
		fmt.Println("Starting Telegram PA Bot")
		internal.StartTelegramPABot(ctx)
	}()

	<-ctx.Done()
	fmt.Println("\nShutdown signal received...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer shutdownCancel()
	
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("error shutting down server: %s\n", err)
	}
	fmt.Println("Graceful shutdown complete.")

}
