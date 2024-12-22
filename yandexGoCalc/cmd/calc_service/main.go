package main

import (
	"fmt"
	"net/http"
	"yandexGoCalc/api"
)

func main() {
	http.HandleFunc("/api/v1/calculate", api.CalculateHandler)
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
