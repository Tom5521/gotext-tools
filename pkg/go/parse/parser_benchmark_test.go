package parse_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/go/parse"
)

func BenchmarkParse(b *testing.B) {
	const input = `package main

import "github.com/leonelquinteros/gotext"

func main(){
	_ = "Hello World"
	a := "Hi world"
	b := "I love onions!"
	
	var eggs string = "sugar"
	
	gotext.Get("MEOW")
	gotext.Get("Hello World")
	gotext.Get(":D")
}`
	bytes := []byte(input)
	parser, err := parse.NewParserFromBytes(bytes, "test.go")
	if err != nil {
		b.Error(err)
	}

	tests := []struct {
		name    string
		options []parse.Option
	}{
		{name: "normal"},
		{name: "extract-all", options: []parse.Option{parse.WithExtractAll(true)}},
	}

	b.ResetTimer()
	for _, t := range tests {
		b.Run(t.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				parser.ParseWithOptions(t.options...)
				b.StopTimer()
				if parser.Error() != nil {
					b.Error(parser.Error())
					b.Skip(parser.Error())
				}
				b.StartTimer()
			}
		})
	}
}
