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

	incidentList := data.NewIncidentsList()

	webSocket := services.NewManager(incidentList)

	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server run")
	})

	//socket
	e.GET("/ws", webSocket.ServerWS)

	e.Logger.Fatal(e.Start(":8080"))
}
