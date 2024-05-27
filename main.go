package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"rando-api/internal"
)

func main() {
	http.HandleFunc("/", internal.GetIndex)
	http.HandleFunc("/page-view", internal.GetPageView)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	fmt.Printf("Running on port %s\n", PORT)
	PORT = fmt.Sprintf(":%s", PORT)

	err := http.ListenAndServe(PORT, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
