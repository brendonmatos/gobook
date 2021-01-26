package components

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/brendonmatos/golive"
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
	interpreter    *interp.Interpreter
	Stdout, Stderr io.Writer
	Blocks         []*Block
}

func NewBook() *golive.LiveComponent {

	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	return golive.NewLiveComponent("Book", &Book{
		Blocks: make([]*Block, 0),
		Stdout: stdout,
		Stderr: stderr,
	})
}

func (b *Book) Mounted(lc *golive.LiveComponent) {
	// create Stdout handler
	b.interpreter = interp.New(interp.Options{
		Stdout: b.Stdout,
		Stderr: b.Stderr,
	})
	b.interpreter.Use(stdlib.Symbols)
	b.interpreter.Use(interp.Exports{
		"gobook": {
			"viewComponent": reflect.ValueOf(b),
			"view":          reflect.ValueOf(lc),
		},
	})

	b.CreateBlock()
}

func (b *Book) CreateBlock() {

	b.Blocks = append(b.Blocks, &Block{})
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
		<div class="bg-gray-300 min-h-screen flex">
			<div class="w-2/3 py-4 mx-auto rounded p-4 font-mono">
				{{ range $index, $block := .Blocks }}
					<div key="{{$index}}" class="bg-white border-gray-200 rounded border p-2 {{ if gt $index 0 }} mt-2 {{ end }}">
						<div class="flex">
							<div class="flex-grow flex flex-col">
								<div class="flex w-full">
									<div class="whitespace-nowrap w-32 text-right pr-2">In [{{$index}}]:</div>
									<textarea placeholder="Code..." 
										class="w-full h-full border" 
										go-live-keydown="ExecuteBlock" 
										go-live-key="Enter" 
										go-live-data-index="{{$index}}" 
										go-live-input="Blocks.{{$index}}.Code"></textarea>
								</div>

								<div class="flex w-full mt-1">
								
									{{ if ne $block.Out "" }}
										<div class="whitespace-nowrap text-right w-32 pr-2">> [{{$index}}]:</div>
										<div class="border w-full bg-black text-white p-2" >
											<span class="font-xm">></span>
											{{ $block.Out }}
										</div>
									{{ end }}
									
									{{ if ne $block.ErrOut "" }}
										<div class="whitespace-nowrap text-right w-32 pr-2">ErrOut [{{$index}}]:</div>
										<div class="border w-full border-red-400">
											{{ $block.ErrOut }}
										</div>
									{{ end }}
									
									
								</div>
							</div>
							<div class="flex pl-2 text-sm">
								<button class="border border-green-500 text-green-500 w-6 h-6" go-live-click="ExecuteBlock" go-live-data-index="{{$index}}">
									▶️
								</button>
								<button class="border border-yellow-500 text-yellow-500 w-6 h-6" go-live-click="DeleteBlock" go-live-data-index="{{$index}}">
									🗑️
								</button>
							</div>
						</div>
					</div>
				{{ end }}
				<button class="w-full rounded border mt-2  bg-white" go-live-click="CreateBlock">➕ Add block</button>
			</div>
			<div class="bg-black w-1/3 p-4 text-white">
				<pre>
				{{ .Stderr }}
				{{ .Stdout }}
				</pre>
			</div>
		</div>
	`
}
