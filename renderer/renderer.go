package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zeozeozeo/imgui"
)

type GetCursorFunc func() (x, y float32)

type Manager struct {
	Filter             ebiten.Filter
	Cache              TextureCache
	Ctx                *imgui.ImGuiContext // ImGui context
	cliptxt            string
	GetCursor          GetCursorFunc
	SyncInputsFn       func()
	SyncCursor         bool
	SyncInputs         bool
	ControlCursorShape bool
	lmask              *ebiten.Image
	ClipMask           bool

	width        float32
	height       float32
	screenWidth  int
	screenHeight int

	inputChars []rune
}

// Text implements ImGui clipboard
func (m *Manager) Text() (string, error) {
	return m.cliptxt, nil
}

// SetText implements ImGui clipboard
func (m *Manager) SetText(text string) {
	m.cliptxt = text
}

// SetDisplaySize sets the display dimensions.
func (m *Manager) SetDisplaySize(width, height float32) {
	m.width = width
	m.height = height
}

// BeginFrame begins a new ImGui frame.
func (m *Manager) BeginFrame() {
	imgui.NewFrame()
}

// EndFrame ends the current ImGui frame.
func (m *Manager) EndFrame() {
	imgui.EndFrame()
}

// New creates a new Manager with a provided font atlas.
func New(fontAtlas *imgui.ImFontAtlas) *Manager {
	imctx := imgui.CreateContext(fontAtlas)
	m := &Manager{
		Cache:              NewCache(),
		Ctx:                imctx,
		SyncCursor:         true,
		SyncInputs:         true,
		ClipMask:           true,
		ControlCursorShape: true,
		inputChars:         make([]rune, 0, 256),
	}

	// build the texture font atlas
	ctx := imgui.GetCurrentContext()

	// GetTexDataAsRGBA32 expects a valid outPixels reference, but all of the
	// other parameters can be nil. we don't need the actual pixels, we just
	// want it to generate the font atlas
	var outPixels []uint32
	ctx.IO.Fonts.GetTexDataAsRGBA32(&outPixels, nil, nil, nil)
	outPixels = nil

	ctx.IO.Fonts.SetTexID(1)
	m.Cache.SetFontAtlasTextureID(1)

	m.setKeyMapping()
	return m
}

// NewWithContext creates a new Manager with a provided ImGui context
func NewWithContext(ctx *imgui.ImGuiContext) *Manager {
	m := &Manager{
		Cache:              NewCache(),
		Ctx:                ctx,
		SyncCursor:         true,
		SyncInputs:         true,
		ClipMask:           true,
		ControlCursorShape: true,
	}
	m.setKeyMapping()
	return m
}

func (m *Manager) controlCursorShape() {
	if !m.ControlCursorShape {
		return
	}

	switch imgui.GetMouseCursor() {
	case imgui.ImGuiMouseCursor_None:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	case imgui.ImGuiMouseCursor_Arrow:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	case imgui.ImGuiMouseCursor_TextInput:
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	case imgui.ImGuiMouseCursor_ResizeAll:
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	case imgui.ImGuiMouseCursor_ResizeEW:
		ebiten.SetCursorShape(ebiten.CursorShapeEWResize)
	case imgui.ImGuiMouseCursor_ResizeNS:
		ebiten.SetCursorShape(ebiten.CursorShapeNSResize)
	case imgui.ImGuiMouseCursor_Hand:
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	default:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
}

func (m *Manager) Update(delta float32) {
	ctx := imgui.GetCurrentContext()
	ctx.IO.DeltaTime = delta

	if m.width > 0 || m.height > 0 {
		ctx.IO.DisplaySize = *imgui.NewImVec2(m.width, m.height)
	} else if m.screenWidth > 0 || m.screenHeight > 0 {
		ctx.IO.DisplaySize = *imgui.NewImVec2(float32(m.screenWidth), float32(m.screenHeight))
	}

	if m.SyncCursor {
		// update cursor position
		if m.GetCursor != nil {
			x, y := m.GetCursor()
			ctx.IO.MousePos = *imgui.NewImVec2(x, y)
		} else {
			mx, my := ebiten.CursorPosition()
			ctx.IO.MousePos = *imgui.NewImVec2(float32(mx), float32(my))
		}

		// update mouse buttons
		ctx.IO.MouseDown[0] = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		ctx.IO.MouseDown[1] = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
		ctx.IO.MouseDown[2] = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)

		// update mouse wheel
		xoff, yoff := ebiten.Wheel()
		ctx.IO.MouseWheel += float32(yoff)
		ctx.IO.MouseWheelH += float32(xoff)
		m.controlCursorShape()
	}

	if m.SyncInputs {
		if m.SyncInputsFn != nil {
			m.SyncInputsFn()
		} else {
			m.inputChars = sendInput(&ctx.IO, m.inputChars)
		}
	}
}

func (m *Manager) Draw(screen *ebiten.Image) {
	m.screenWidth = screen.Bounds().Dx()
	m.screenHeight = screen.Bounds().Dy()
	imgui.Render()

	if m.ClipMask {
		if m.lmask == nil {
			w, h := screen.Size()
			m.lmask = ebiten.NewImage(w, h)
		} else {
			w1, h1 := screen.Size()
			w2, h2 := m.lmask.Size()
			if w1 != w2 || h1 != h2 {
				m.lmask.Dispose()
				m.lmask = ebiten.NewImage(w1, h1)
			}
		}
		RenderMasked(screen, m.lmask, imgui.GetDrawData(), m.Cache, m.Filter)
	} else {
		Render(screen, imgui.GetDrawData(), m.Cache, m.Filter)
	}
}
