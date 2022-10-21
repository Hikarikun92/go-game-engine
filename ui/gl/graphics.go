package gl

import (
	"github.com/Hikarikun92/go-game-engine/ui"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type graphicsImpl struct {
	shaderProgram uint32
}

func (g *graphicsImpl) DrawImage(image ui.Image, x int, y int) {
	img := image.(imageImpl)

	model := mgl32.Translate3D(float32(x), float32(y), -1)
	model = model.Mul4(mgl32.Scale3D(img.width, img.height, 1.0))

	modelUniform := gl.GetUniformLocation(g.shaderProgram, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, img.textureId)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
