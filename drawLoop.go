package oak

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
	"golang.org/x/exp/shiny/screen"
)

var (
	imageBlack = image.Black
	// DrawTicker is an unused parallel to LogicTicker to set the draw framerate
	DrawTicker *timing.DynamicTicker
)

// DrawLoop
// Unless told to stop, the draw channel will repeatedly
// 1. draw black to a temporary buffer
// 2. draw all elements onto the temporary buffer.
// 3. scale the buffer's data at the viewport's position to a texture.
// 4. publish the texture to display on screen.
func drawLoop() {
	<-drawCh

	setPanicOnFault()

	tx, err := screenControl.NewTexture(winBuffer.Bounds().Max)
	if err != nil {
		panic(err)
	}

	draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, screen.Src)
	drawLoopPublish(tx)
	fmt.Println("First publish done")

	DrawTicker = timing.NewDynamicTicker()
	fmt.Println("Dynamic Ticker gotten")
	DrawTicker.SetTick(timing.FPSToDuration(DrawFrameRate))

	dlog.Verb("Draw Loop Start")
	fmt.Println("Draw Loop Start")
	for {
	drawSelect:
		select {
		case <-windowUpdateCh:
			<-windowUpdateCh
		case <-drawCh:
			dlog.Verb("Got something from draw channel")
			<-drawCh
			dlog.Verb("Starting loading")
			fmt.Println("Starting loading")
			for {
				<-DrawTicker.C
				draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, screen.Src)
				if LoadingR != nil {
					LoadingR.Draw(winBuffer.RGBA())
				}
				drawLoopPublish(tx)

				select {
				case <-drawCh:
					break drawSelect
				case viewPoint := <-viewportCh:
					dlog.Verb("Got something from viewport channel (waiting on draw)")
					updateScreen(viewPoint[0], viewPoint[1])
				default:
				}
			}
		case viewPoint := <-viewportCh:
			dlog.Verb("Got something from viewport channel")
			updateScreen(viewPoint[0], viewPoint[1])
		case <-DrawTicker.C:
			draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, screen.Src)
			render.PreDraw()
			render.GlobalDrawStack.Draw(winBuffer.RGBA(), ViewPos, ScreenWidth, ScreenHeight)
			drawLoopPublish(tx)
		}
	}
}

var (
	drawLoopPublish = drawLoopPublishDef
)
