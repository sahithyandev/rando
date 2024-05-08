package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello from go")
	// internal.GetPageView()
	// http.HandleFunc("/", internal.GetPageView)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	PORT = fmt.Sprintf(":%s", PORT)

	err := http.ListenAndServe(PORT)
}
