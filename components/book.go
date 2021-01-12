package components

import (
	"fmt"
	"strconv"

	"github.com/brendonferreira/golive"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Block struct {
	Code   string
	Out    string
	ErrOut string
}

type Book struct {
	golive.LiveComponentWrapper
	interpreter *interp.Interpreter
	Blocks      []*Block
}

func NewBook() *golive.LiveComponent {
	return golive.NewLiveComponent("Book", &Book{
		Blocks: make([]*Block, 0),
	})
}

func (b *Book) Mounted(_ *golive.LiveComponent) {
	// create Stdout handler
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
	block := b.Blocks[index]

	// Clear block
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
	return `
		<div class="bg-gray-300 py-4 min-h-screen">
			<div class="container bg-white mx-auto rounded p-4 font-mono" style="max-width: 800px">
				{{ range $index, $block := .Blocks }}
					<div key="{{$index}}" class="border-gray-200 rounded border p-2 {{ if gt $index 0 }} mt-2 {{ end }}">
						<div class="flex justify-between pb-2">
							<button class="border border-green-500 text-green-500 px-4" go-live-click="ExecuteBlock" go-live-data-index="{{$index}}">
								▶️
							</button>
							<button class="border border-yellow-500 text-yellow-500 px-4" go-live-click="DeleteBlock" go-live-data-index="{{$index}}">
								🗑️
							</button>
						</div>
						<div class="">
							<textarea placeholder="Code..." go-live-keydown-enter="ExecuteBlock" go-live-data-index="{{$index}}" class="w-full h-full border" go-live-input="Blocks.{{$index}}.Code"></textarea>
						</div>
						{{ if ne $block.Out "" }}
							<div class="border bg-black text-white p-4" >
								<span class="font-xm">></span>
								{{ $block.Out }}
							</div>
						{{ end }}
						{{ if ne $block.ErrOut "" }}
							<div class="border border-red-400">
								{{ $block.ErrOut }}
							</div>
						{{ end }}
					</div>
				{{ end }}
				<button class="w-full rounded border mt-2" go-live-click="CreateBlock">➕ Add block</button>
			</div>
		</div>
	`
}
