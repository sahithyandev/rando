package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"rando-api/internal"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/page-view", internal.GetPageView)
	mux.HandleFunc("/live-users-count", internal.GetLiveUsersCount)

	internal.StartTelegramPABot()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	fmt.Printf("Running on port %s\n", PORT)
	PORT = fmt.Sprintf(":%s", PORT)

	handler := cors.Default().Handler(mux)
	err = http.ListenAndServe(PORT, handler)
	fmt.Println("Server is running.")

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
