package gl

import (
	"github.com/Hikarikun92/go-game-engine/key"
	"github.com/Hikarikun92/go-game-engine/ui"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
)

type glWindowManager struct {
}

func NewWindowManager() ui.WindowManager {
	return &glWindowManager{}
}

type windowImpl struct {
	glfwWindow          *glfw.Window
	vertexArrayObject   uint32
	vertexBufferObject  uint32
	elementBufferObject uint32
	shaderProgram       uint32
}

/*
References:
https://learnopengl.com/Getting-started/Textures
https://learnopengl.com/In-Practice/2D-Game/Rendering-Sprites
https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
*/

func (*glWindowManager) CreateMainWindow() ui.Window {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize GLFW:", err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Example game", nil, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	shaderProgram, err := newShaderProgram()
	if err != nil {
		log.Fatalln("Failed to create shader program:", err)
	}

	gl.UseProgram(shaderProgram)

	textureUniform := gl.GetUniformLocation(shaderProgram, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	vertices := []float32{
		// pos    // tex
		0.0, 1.0, 0.0, 0.0, //bottom left
		1.0, 0.0, 1.0, 1.0, //top right
		0.0, 0.0, 0.0, 1.0, //top left
		1.0, 1.0, 1.0, 0.0, //bottom right
	}
	indices := []int32{
		0, 1, 2,
		0, 3, 1,
	}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	//position and texture attributes (the "location = 0" in the shader)
	gl.VertexAttribPointerWithOffset(0, 4, gl.FLOAT, false, 4*4 /* 4 values per vertex * 4 bytes per value */, 0)
	gl.EnableVertexAttribArray(0)

	window.Show()

	return &windowImpl{
		glfwWindow:          window,
		vertexArrayObject:   vao,
		vertexBufferObject:  vbo,
		elementBufferObject: ebo,
		shaderProgram:       shaderProgram,
	}
}

func (w *windowImpl) SetKeyListener(listener key.Listener) {
	w.glfwWindow.SetKeyCallback(func(w *glfw.Window, glfwKey glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		k := translateKey(glfwKey)
		if k == key.UNKNOWN {
			return
		}

		if action == glfw.Press {
			listener.KeyPressed(k)
		} else if action == glfw.Release {
			listener.KeyReleased(k)
		}
	})
}

func (w *windowImpl) CreateImageLoader() ui.ImageLoader {
	return &imageLoaderImpl{}
}

func (w *windowImpl) CreateGraphics() ui.Graphics {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	return &graphicsImpl{}
}

func (w *windowImpl) ShouldClose() bool {
	return w.glfwWindow.ShouldClose()
}

func (w *windowImpl) Update() {
	w.glfwWindow.SwapBuffers()
	glfw.PollEvents()
}

func (w *windowImpl) Destroy() {
	gl.DeleteVertexArrays(1, &w.vertexArrayObject)
	gl.DeleteBuffers(1, &w.vertexBufferObject)
	gl.DeleteBuffers(1, &w.elementBufferObject)
	gl.DeleteProgram(w.shaderProgram)

	w.glfwWindow.Destroy()
	glfw.Terminate()
}
