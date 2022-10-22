package gl

import (
	"github.com/Hikarikun92/go-game-engine/ui"
	"github.com/go-gl/gl/v4.1-core/gl"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

type imageLoaderImpl struct {
}

type imageImpl struct {
	textureId uint32
	width     float32
	height    float32
}

func (i *imageLoaderImpl) LoadImage(file string) ui.Image {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v", file, err)
	}

	//Decode the image to a know structure (using the imports with _)
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	rgbaSize := rgba.Rect.Size()

	if rgba.Stride != rgbaSize.X*4 {
		log.Fatalln("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)

	//Create an OpenGL texture with the image data
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgbaSize.X), int32(rgbaSize.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return imageImpl{
		textureId: texture,
		width:     float32(rgbaSize.X),
		height:    float32(rgbaSize.Y),
	}
}

func (i *imageLoaderImpl) UnloadImage(image ui.Image) {
	img := image.(imageImpl)
	gl.DeleteTextures(1, &img.textureId)
}
