package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/srq/signalflux/internal/api"
	"github.com/srq/signalflux/internal/config"
)


func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(" No env file found")
	}
}
func main() {
	cfg := config.LoadConfig()

	fmt.Println("Server running on port:", cfg.ServerPort)

	router := api.NewRouter()

	http.ListenAndServe(":8080", router)

}
