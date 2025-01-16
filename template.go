//go:build ignore
// +build ignore

package main

import po "github.com/leonelquinteros/gotext"

// This should cover almost all get methods.
func main() {
	po.Get("Hi %s!\n", "stranger")
	po.Get("Hi %s!", "stranger")
	po.Get("Hi %s!", "stranger")
	po.Get(`HIIII %s

"Hello World"

\n

a`, "stranger")
	po.GetC("Hi %s!", "formal", "stranger") // Hello stranger!
	po.GetC("Hi %s!", "formal", "stranger")
	po.GetC("Hi %s!", "casual", "stranger")
	po.GetC("Hi %s!", "casual", "stranger")
	po.GetD("default", "Hello World! %d", 1234)
	po.GetN("I want %d apple", "I want %d apples", 1, 3)
	po.GetND("default", "Hello World!", "Hello Worlds!", 1)
	po.GetNC("Hello World!", "Hello Worlds!", 1, "mars")
	po.GetNDC("moon", "Hi stranger, I'm %s", "Hi strangers, I'm %s", 1, "mars", "Tom!")
}
