package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/zeozeozeo/imgui"
)

var keys = map[imgui.ImGuiKey]int{
	imgui.ImGuiKey_Tab:        int(ebiten.KeyTab),
	imgui.ImGuiKey_LeftArrow:  int(ebiten.KeyLeft),
	imgui.ImGuiKey_RightArrow: int(ebiten.KeyRight),
	imgui.ImGuiKey_UpArrow:    int(ebiten.KeyUp),
	imgui.ImGuiKey_DownArrow:  int(ebiten.KeyDown),
	imgui.ImGuiKey_PageUp:     int(ebiten.KeyPageUp),
	imgui.ImGuiKey_PageDown:   int(ebiten.KeyPageDown),
	imgui.ImGuiKey_Home:       int(ebiten.KeyHome),
	imgui.ImGuiKey_End:        int(ebiten.KeyEnd),
	imgui.ImGuiKey_Insert:     int(ebiten.KeyInsert),
	imgui.ImGuiKey_Delete:     int(ebiten.KeyDelete),
	imgui.ImGuiKey_Backspace:  int(ebiten.KeyBackspace),
	imgui.ImGuiKey_Space:      int(ebiten.KeySpace),
	imgui.ImGuiKey_Enter:      int(ebiten.KeyEnter),
	imgui.ImGuiKey_Escape:     int(ebiten.KeyEscape),
	imgui.ImGuiKey_A:          int(ebiten.KeyA),
	imgui.ImGuiKey_C:          int(ebiten.KeyC),
	imgui.ImGuiKey_V:          int(ebiten.KeyV),
	imgui.ImGuiKey_X:          int(ebiten.KeyX),
	imgui.ImGuiKey_Y:          int(ebiten.KeyY),
	imgui.ImGuiKey_Z:          int(ebiten.KeyZ),
}

func sendInput(io *imgui.ImGuiIO, inputChars []rune) []rune {
	io.KeyAlt = ebiten.IsKeyPressed(ebiten.KeyAlt)
	io.KeyShift = ebiten.IsKeyPressed(ebiten.KeyShift)
	io.KeyCtrl = ebiten.IsKeyPressed(ebiten.KeyControl)
	// TODO: KeySuper

	inputChars = ebiten.AppendInputChars(inputChars)
	if len(inputChars) > 0 {
		io.AddInputCharacters(string(inputChars))
		inputChars = inputChars[:0]
	}
	for _, key := range keys {
		if inpututil.IsKeyJustPressed(ebiten.Key(key)) {
			io.KeysDown[key] = true
		}
		if inpututil.IsKeyJustReleased(ebiten.Key(key)) {
			io.KeysDown[key] = false
		}
	}
	return inputChars
}

func (m *Manager) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	ctx := imgui.GetCurrentContext()
	for imguiKey, nativeKey := range keys {
		ctx.IO.KeyMap[imguiKey] = int32(nativeKey)
	}
}