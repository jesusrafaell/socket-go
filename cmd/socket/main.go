package main

import (
	"crashsaver/websocket/data"
	"crashsaver/websocket/pkg/socket"
	"net/http"

	// "github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %s", err)
	// }

	// // Getting and using a value from .env
	// app := os.Getenv("APP")

	// log.Println(app)

	NewServer()
}

func NewServer() {

	incidentList := data.NewIncidents()

	webSocket := socket.NewManager(incidentList)

	e := echo.New()

	//initn db

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server run")
	})

	// api = e.Group("/api", handlers.WithAuthentication)

	//socket
	e.GET("/ws", webSocket.ServerWS)

	e.Logger.Fatal(e.Start(":8080"))
}
