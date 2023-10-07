package main

import (
	socketService "crashsaver/websocket/services/socket"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	newServer()
}

func newServer() {
	e := echo.New()

	webSocket := socketService.NewSocket()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server run")
	})

	e.GET("/ws", webSocket.ServerWS)

	e.Logger.Fatal(e.Start(":8000"))
}
