package renderer

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zeozeozeo/imgui"
)

func getTexture(width, height int32, pixels []uint32) *ebiten.Image {
	n := width * height
	pix := make([]uint8, n*4)

	// NOTE: ebiten expects premultiplied-alpha pixels, but imgui
	// exports bitmaps in straight-alpha.
	for i := int32(0); i < n; i++ {
		r := pixels[i] & 0xFF
		g := pixels[i] >> 8 & 0xFF
		b := pixels[i] >> 16 & 0xFF
		a := pixels[i] >> 24 & 0xFF

		alpha := uint16(a)
		pix[4*i] = uint8((uint16(r)*alpha + 127) / 255)
		pix[4*i+1] = uint8((uint16(g)*alpha + 127) / 255)
		pix[4*i+2] = uint8((uint16(b)*alpha + 127) / 255)
		pix[4*i+3] = uint8(alpha)
	}

	img := ebiten.NewImage(int(width), int(height))
	img.WritePixels(pix)
	return img
}

func convertVertices(verts []imgui.ImDrawVert) []ebiten.Vertex {
	vertices := make([]ebiten.Vertex, len(verts))

	for i := 0; i < len(verts); i++ {
		// ImGui uses 32 bit unsigned integers to store colors,
		// but Ebiten uses 0-1 floating point numbers
		vertices[i] = ebiten.Vertex{
			SrcX:   verts[i].Uv.X(),
			SrcY:   verts[i].Uv.Y(),
			DstX:   verts[i].Pos.X(),
			DstY:   verts[i].Pos.Y(),
			ColorR: float32(verts[i].Col&0xFF) / 255,
			ColorG: float32(verts[i].Col>>8&0xFF) / 255,
			ColorB: float32(verts[i].Col>>16&0xFF) / 255,
			ColorA: float32(verts[i].Col>>24&0xFF) / 255,
		}
	}
	return vertices
}

func vcopy(v []ebiten.Vertex) []ebiten.Vertex {
	cl := make([]ebiten.Vertex, len(v))
	copy(cl, v)
	return cl
}

// Render the ImGui drawData into the target *ebiten.Image
func Render(target *ebiten.Image, drawData *imgui.ImDrawData, txcache TextureCache, dfilter ebiten.Filter) {
	render(target, nil, drawData, txcache, dfilter)
}

// RenderMasked renders the ImGui drawData into the target *ebiten.Image with ebiten.CompositeModeCopy for masking
func RenderMasked(target *ebiten.Image, mask *ebiten.Image, drawData *imgui.ImDrawData, txcache TextureCache, dfilter ebiten.Filter) {
	render(target, mask, drawData, txcache, dfilter)
}

func lerp(a, b int, t float32) float32 {
	return float32(a)*(1-t) + float32(b)*t
}

func vmultiply(v, vbuf []ebiten.Vertex, bmin, bmax image.Point) {
	for i := range vbuf {
		vbuf[i].SrcX = lerp(bmin.X, bmax.X, v[i].SrcX)
		vbuf[i].SrcY = lerp(bmin.Y, bmax.Y, v[i].SrcY)
	}
}

func indicesToUint16(input []imgui.ImDrawIdx) []uint16 {
	out := make([]uint16, len(input))
	for i, val := range input {
		out[i] = uint16(val)
	}
	return out
}

func render(target *ebiten.Image, mask *ebiten.Image, drawData *imgui.ImDrawData, txcache TextureCache, dfilter ebiten.Filter) {
	if !drawData.Valid {
		return
	}
	targetw, targeth := target.Size()

	// image and triangle options
	opt := &ebiten.DrawTrianglesOptions{
		Filter: dfilter,
	}

	var opt2 *ebiten.DrawImageOptions
	if mask != nil {
		opt2 = &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
		}
	}

	for _, clist := range drawData.CmdLists {
		// var indexBufferOffset uint32
		vertices := convertVertices(clist.VtxBuffer)
		verticesMul := vcopy(vertices)

		// draw command buffer
		for _, cmd := range clist.CmdBuffer {
			if cmd.UserCallback != nil {
				cmd.UserCallback(clist, &cmd)
				continue
			}
			// has no user callbacks

			clipRect := cmd.ClipRect
			texture := txcache.GetTexture(cmd.GetTexID())
			vmultiply(vertices, verticesMul, texture.Bounds().Min, texture.Bounds().Max)

			indices := indicesToUint16(clist.IdxBuffer[cmd.IdxOffset : cmd.IdxOffset+cmd.ElemCount])

			// if has no mask, draw triangles like normal
			if mask == nil || (clipRect.X() == 0 && clipRect.Y() == 0 && clipRect.Z() == float32(targetw) && clipRect.W() == float32(targeth)) {
				target.DrawTriangles(verticesMul, indices, texture, opt)
				continue
			}

			// has a clip mask
			mask.Clear()
			opt2.GeoM.Reset()
			opt2.GeoM.Translate(float64(clipRect.X()), float64(clipRect.Y()))

			mask.DrawTriangles(verticesMul, indices, texture, opt)

			target.DrawImage(mask.SubImage(image.Rectangle{
				Min: image.Pt(int(clipRect.X()), int(clipRect.Y())),
				Max: image.Pt(int(clipRect.Z()), int(clipRect.W())),
			}).(*ebiten.Image), opt2)
		}
	}

}

/*
func render(target *ebiten.Image, mask *ebiten.Image, drawData *imgui.ImDrawData, txcache TextureCache, dfilter ebiten.Filter) {
	targetw, targeth := target.Size()
	if !drawData.Valid {
		return
	}

	opt := &ebiten.DrawTrianglesOptions{
		Filter: dfilter,
	}
	var opt2 *ebiten.DrawImageOptions
	if mask != nil {
		opt2 = &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
		}
	}

	for _, clist := range drawData.CmdLists {
		var indexBufferOffset uint32
		indexBuffer := clist.IdxBuffer

		vertices := getVertices(clist.VtxBuffer)
		vbuf := vcopy(vertices)

		for _, cmd := range clist.CmdBuffer {
			ecount := cmd.ElemCount

			if cmd.UserCallback != nil {
				cmd.UserCallback(clist, &cmd)
			} else {
				clipRect := cmd.ClipRect
				texid := cmd.GetTexID() // imgui.ImTextureID
				tx := txcache.GetTexture(texid)
				vmultiply(vertices, vbuf, tx.Bounds().Min, tx.Bounds().Max)

				if mask == nil || (clipRect.X() == 0 && clipRect.Y() == 0 && clipRect.Z() == float32(targetw) && clipRect.W() == float32(targeth)) {
					target.DrawTriangles(
						vbuf,
						indicesToUint16(indexBuffer[indexBufferOffset:indexBufferOffset+ecount]),
						tx,
						opt,
					)
				} else {
					mask.Clear()
					opt2.GeoM.Reset()
					opt2.GeoM.Translate(float64(clipRect.X()), float64(clipRect.Y()))

					mask.DrawTriangles(
						vbuf,
						indicesToUint16(indexBuffer[indexBufferOffset:indexBufferOffset+ecount]),
						tx,
						opt,
					)

					target.DrawImage(mask.SubImage(image.Rectangle{
						Min: image.Pt(int(clipRect.X()), int(clipRect.Y())),
						Max: image.Pt(int(clipRect.Z()), int(clipRect.W())),
					}).(*ebiten.Image), opt2)
				}
			}
			indexBufferOffset += ecount
		}
	}
}
*/
