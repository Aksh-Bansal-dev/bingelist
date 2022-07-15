package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/Aksh-Bansal-dev/bingelist/pkg/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("ENV") == "dev" {
		log.SetFlags(log.Lshortfile)
	}
	routes.Routes()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}
	log.Println("Server started at port:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
