package main

import (
	"fmt"
	"github.com/brendonmatos/gobook/components"
	"github.com/brendonmatos/golive"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/skratchdot/open-golang/open"
	"net"
)

func main() {

	app := fiber.New()
	liveServer := golive.NewServer()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	app.Get("/", liveServer.CreateHTMLHandler(components.NewBook, golive.PageContent{
		Lang:  "us",
		Title: "Hello world",
		Head:  `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/tailwindcss/2.0.2/tailwind.min.css" integrity="sha512-+WF6UMXHki/uCy0vATJzyA9EmAcohIQuwpNz0qEO+5UeE5ibPejMRdFuARSrl1trs3skqie0rY/gNiolfaef5w==" crossorigin="anonymous" />`,
	}))

	app.Get("/ws", websocket.New(liveServer.HandleWSRequest))

	_ = open.Start("http://" + listener.Addr().String())

	panic(app.Listener(listener))
}
