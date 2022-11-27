package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/zeozeozeo/ebitengine-imgui-go/renderer"
	"github.com/zeozeozeo/imgui"
)

func main() {
	// create a new imgui context, you can provide your font atlas if you want,
	// but we'll just pass nil for the default one
	mgr := renderer.New(nil)
	mgr.Ctx.IO.IniFilename = "" // disable imgui.ini

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	gg := &Game{
		mgr:            mgr,
		deviceScale:    ebiten.DeviceScaleFactor(),
		demoWindowOpen: true,
	}

	ebiten.RunGame(gg)
}

type Game struct {
	mgr *renderer.Manager

	deviceScale    float64
	retina         bool
	w, h           int
	demoWindowOpen bool
	consoleOutput  []string
}

func (game *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"tps: %.2f\nfps: %.2f\n(c)lipmask: %t",
			ebiten.ActualTPS(),
			ebiten.ActualFPS(),
			game.mgr.ClipMask,
		),
		10,
		2,
	)

	game.mgr.Draw(screen)
}

func (game *Game) Update() error {
	// update imgui state
	game.mgr.Update(1.0 / float32(ebiten.ActualTPS()))

	// disable or enable the clipmask when C is pressed
	// usually you'd always have this enabled
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		game.mgr.ClipMask = !game.mgr.ClipMask
	}

	game.mgr.BeginFrame() // start a new imgui frame

	// test window
	imgui.Begin("Test window", nil, 0) // 0 = default window flags
	{
		imgui.Text("Hello, world!")
		if imgui.Button("button") {
			fmt.Println("button pressed")
			game.consoleOutput = append(game.consoleOutput, "button pressed")
		}

		// draw a second button on the same line
		imgui.SameLine(0, 4)
		if imgui.Button("button 2") {
			fmt.Println("button 2 pressed")
			game.consoleOutput = append(game.consoleOutput, "button 2 pressed")
		}
	}
	imgui.End()

	// options window
	imgui.SetNextWindowPos(imgui.NewImVec2(300, 300), imgui.ImGuiCond_FirstUseEver, *imgui.NewImVec2(0, 0))
	imgui.Begin("Options", nil, 0)
	{
		imgui.Checkbox("Retina", &game.retina)
		imgui.Checkbox("Demo window", &game.demoWindowOpen)
	}
	imgui.End()

	// console window
	imgui.SetNextWindowPos(imgui.NewImVec2(50, 200), imgui.ImGuiCond_FirstUseEver, *imgui.NewImVec2(0, 0))
	imgui.SetNextWindowSize(imgui.NewImVec2(200, 300), imgui.ImGuiCond_FirstUseEver)
	imgui.Begin("Console", nil, 0)
	{
		imgui.Text("console output:")
		for idx, s := range game.consoleOutput {
			imgui.PushID(int32(idx))
			imgui.Text("%s", s)
			imgui.PopID()
		}
	}
	imgui.End()

	// demo window (WARNING: this crashes a lot :p)
	if game.demoWindowOpen {
		imgui.ShowDemoWindow(&game.demoWindowOpen)
	}

	game.mgr.EndFrame()
	return nil
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if game.retina {
		m := ebiten.DeviceScaleFactor()
		game.w = int(float64(outsideWidth) * m)
		game.h = int(float64(outsideHeight) * m)
	} else {
		game.w = outsideWidth
		game.h = outsideHeight
	}

	game.mgr.SetDisplaySize(float32(game.w), float32(game.h))
	return game.w, game.h
}
