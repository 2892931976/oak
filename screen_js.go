//+build js

package oak

import (
	"fmt"
	"image"

	"github.com/gopherjs/gopherjs/js"
	"golang.org/x/exp/shiny/screen"
)

type JSScreen struct{}

func (jss *JSScreen) NewBuffer(p image.Point) (screen.Buffer, error) {
	fmt.Println("New JS Buffer")
	rect := image.Rect(0, 0, p.X, p.Y)
	rgba := image.NewRGBA(rect)
	buffer := &JSBuffer{
		rect,
		rgba,
	}
	return buffer, nil
}
func (jss *JSScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	fmt.Println("New JS Window")
	jsc := new(JSWindow)

	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Get("style").Set("display", "block")
	canvas.Set("width", ScreenWidth)
	canvas.Set("height", ScreenHeight)
	jsc.ctx = canvas.Call("getContext", "2d")
	bdy := document.Get("body")
	bdy.Call("appendChild", canvas)

	return jsc, nil
}

func (jss *JSScreen) NewTexture(p image.Point) (screen.Texture, error) {
	fmt.Println("New JS Texture")
	txt := new(JSTexture)
	return txt, nil
}
