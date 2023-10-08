package main

import (
	"crashsaver/websocket/data"
	"crashsaver/websocket/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	NewServer()
}

func NewServer() {

	// Crear una instancia única de la lista de incidentes
	incidentList := data.NewIncidentsList()

	// Crear el servidor WebSocket
	webSocket := services.NewManager(incidentList)

	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server run")
	})

	//scoket
	e.GET("/ws", webSocket.ServerWS)

	e.Logger.Fatal(e.Start(":8000"))
}
