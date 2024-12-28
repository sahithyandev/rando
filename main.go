package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"rando-api/internal"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/page-view", internal.GetPageView)
	mux.HandleFunc("/live-users-count", internal.GetLiveUsersCount)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	fmt.Printf("Running on port %s\n", PORT)
	PORT = fmt.Sprintf(":%s", PORT)

	handler := cors.Default().Handler(mux)
	err := http.ListenAndServe(PORT, handler)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
