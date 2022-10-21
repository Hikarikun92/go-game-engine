package gl

import (
	"github.com/Hikarikun92/go-game-engine/ui"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type graphicsImpl struct {
}

func (g *graphicsImpl) DrawImage(image ui.Image, x int, y int) {
	img := image.(imageImpl)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, img.textureId)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
