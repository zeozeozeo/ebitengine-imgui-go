package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zeozeozeo/imgui"
)

type TextureCache interface {
	FontAtlasTextureID() imgui.ImTextureID
	SetFontAtlasTextureID(id imgui.ImTextureID)
	GetTexture(id imgui.ImTextureID) *ebiten.Image
	SetTexture(id imgui.ImTextureID, img *ebiten.Image)
	RemoveTexture(id imgui.ImTextureID)
	ResetFontAtlasCache(filter ebiten.Filter)
}

type textureCache struct {
	fontAtlasID    imgui.ImTextureID
	fontAtlasImage *ebiten.Image
	cache          map[imgui.ImTextureID]*ebiten.Image
	dfilter        ebiten.Filter
}

var _ TextureCache = (*textureCache)(nil)

func (c *textureCache) getFontAtlas() *ebiten.Image {
	if c.fontAtlasImage == nil {
		ctx := imgui.GetCurrentContext()
		var outPixels []uint32
		var outWidth int32
		var outHeight int32
		ctx.IO.Fonts.GetTexDataAsRGBA32(&outPixels, &outWidth, &outHeight, nil)
		c.fontAtlasImage = getTexture(outWidth, outHeight, outPixels)
	}

	return c.fontAtlasImage
}

func (c *textureCache) FontAtlasTextureID() imgui.ImTextureID {
	return c.fontAtlasID
}

func (c *textureCache) SetFontAtlasTextureID(id imgui.ImTextureID) {
	c.fontAtlasID = id
	// c.fontAtlasImage = nil
}

func (c *textureCache) GetTexture(id imgui.ImTextureID) *ebiten.Image {
	if id != c.fontAtlasID {
		if im, ok := c.cache[id]; ok {
			return im
		}
	}
	return c.getFontAtlas()
}

func (c *textureCache) SetTexture(id imgui.ImTextureID, img *ebiten.Image) {
	c.cache[id] = img
}

func (c *textureCache) RemoveTexture(id imgui.ImTextureID) {
	delete(c.cache, id)
}

func (c *textureCache) ResetFontAtlasCache(filter ebiten.Filter) {
	c.fontAtlasImage = nil
	c.dfilter = filter
}

func NewCache() TextureCache {
	return &textureCache{
		fontAtlasID:    1,
		cache:          make(map[imgui.ImTextureID]*ebiten.Image),
		fontAtlasImage: nil,
	}
}
