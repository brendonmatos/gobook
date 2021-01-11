package main

import (
	"fmt"
	"github.com/brendonferreira/golive"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io/ioutil"
	"strconv"
)


type Block struct {
	Code string
	Out  string
	ErrOut string
}

type Book struct {
	golive.LiveComponentWrapper
	interpreter *interp.Interpreter
	Blocks []*Block
}

func NewBook() *golive.LiveComponent {
	fmt.Println()
	return golive.NewLiveComponent("Book", &Book{
		Blocks:               make([]*Block, 0),
	})
}

func (b *Book) Mounted(_ *golive.LiveComponent) {
	b.interpreter = interp.New(interp.Options{})
	b.interpreter.Use(stdlib.Symbols)

	b.CreateBlock()
}

func (b *Book) CreateBlock() {
	b.Blocks = append(b.Blocks, &Block{
		Code: "",
		Out:  "",
	})
}


func (b *Book) DeleteBlock(data map[string]string) {
	index, _ := strconv.ParseInt(data["index"], 10, 8)

	b.Blocks = append(b.Blocks[:index], b.Blocks[index+1:]...)
}


func (b *Book) ExecuteBlock(data map[string]string) {
	index, _ := strconv.ParseInt(data["index"], 10, 8)
	block := b.Blocks[ index ]

	block.Out = ""
	block.ErrOut = ""
	output, err := b.interpreter.Eval(block.Code)
	if err != nil {
		block.ErrOut = fmt.Sprintf("go interpreter error: %s", err)
		return
	}

	block.Out = fmt.Sprintf("%v", output)
}


func (b *Book) TemplateHandler(_ *golive.LiveComponent) string {
	a, _ := ioutil.ReadFile("./cmd/book.gohtml")
	return string(a)
}

func main() {

	app := fiber.New()
	liveServer := golive.NewServer()

	app.Get("/", liveServer.CreateHTMLHandler(NewBook, golive.PageContent{
		Lang:  "us",
		Title: "Hello world",
		Head: `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/tailwindcss/2.0.2/tailwind.min.css" integrity="sha512-+WF6UMXHki/uCy0vATJzyA9EmAcohIQuwpNz0qEO+5UeE5ibPejMRdFuARSrl1trs3skqie0rY/gNiolfaef5w==" crossorigin="anonymous" />`,
	}))

	app.Get("/ws", websocket.New(liveServer.HandleWSRequest))

	_ = app.Listen(":3000")
}

