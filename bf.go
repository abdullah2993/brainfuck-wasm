//go:build wasm

package main

import (
	"log"
	"syscall/js"
)

var inputFn js.Value
var outputFn js.Value

func getInput() byte {
	if inputFn.IsUndefined() {
		log.Fatal("no input handler found")
	}

	res := inputFn.Invoke()
	return res.String()[0]
}

func writeOutput(b byte) {
	if outputFn.IsUndefined() {
		log.Fatal("no output handler found")
	}

	outputFn.Invoke(string(b))
}

func RunBF(v js.Value) {
	runBF(v.String())
}

func runBF(source string) {
	data := make([]byte, 65535)
	loopStack := []int{}
	var pc uint16

	for i := 0; i < len(source); i++ {
		switch source[i] {
		case '+':
			data[pc]++
		case '-':
			data[pc]--
		case '>':
			pc++
		case '<':
			pc--
		case '[':
			if data[pc] == 0 {
				loopCount := 0
				for i < len(source) {
					if source[i] == '[' {
						loopCount++
					} else if source[i] == ']' {
						loopCount--
						if loopCount == 0 {
							break
						}
					}
					i++
					continue
				}
			} else {
				loopStack = append(loopStack, i)
			}
		case ']':
			if data[pc] == 0 {
				loopStack = loopStack[0 : len(loopStack)-1]
			} else {
				i = loopStack[len(loopStack)-1]
			}
		case '.':
			writeOutput(data[pc])
		case ',':
			data[pc] = getInput()
		default:
			continue
		}
	}
}

func main() {
	log.Println("setting up wasm")
	js.Global().Set("registerInputHandler", js.FuncOf(func(this js.Value, args []js.Value) any {
		inputFn = args[0]
		return js.Undefined()
	}))

	js.Global().Set("registerOutputHandler", js.FuncOf(func(this js.Value, args []js.Value) any {
		outputFn = args[0]
		return js.Undefined()
	}))

	js.Global().Set("runBF", js.FuncOf(func(this js.Value, args []js.Value) any {
		RunBF(args[0])
		return js.Undefined()
	}))

	select {}
}
