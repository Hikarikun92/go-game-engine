package gl

import (
	"github.com/Hikarikun92/go-game-engine/cursor"
	"github.com/Hikarikun92/go-game-engine/key"
	"github.com/Hikarikun92/go-game-engine/settings"
	"github.com/Hikarikun92/go-game-engine/ui"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
func (*glWindowManager) CreateMainWindow(settings *settings.Settings) ui.Window {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	//Initialize GLFW and create an invisible window
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize GLFW:", err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(settings.Width, settings.Height, settings.WindowTitle, nil, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}
	window.MakeContextCurrent()

	//Center the window on screen
	videoMode := glfw.GetPrimaryMonitor().GetVideoMode()
	windowX := (videoMode.Width - settings.Width) / 2
	windowY := (videoMode.Height - settings.Height) / 2

	window.SetPos(windowX, windowY)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	//Create the main shader program
	shaderProgram, err := newShaderProgram()
	if err != nil {
		log.Fatalln("Failed to create shader program:", err)
	}

	gl.UseProgram(shaderProgram)

	//Usually you would set HEIGHT as the bottom value and 0 as the top, but I'm deliberately inverting it here
	projection := mgl32.Ortho2D(0.0, float32(settings.Width), 0, float32(settings.Height))

	//Set the projection to the shader
	projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	//Set the sampler for textures
	textureUniform := gl.GetUniformLocation(shaderProgram, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	//Create a square that will be used to draw the images
	vertices := []float32{
		// pos    // tex
		0.0, 1.0, 0.0, 0.0, //top left vertex, bottom left texture
		1.0, 0.0, 1.0, 1.0, //bottom right vertex, top right texture
		0.0, 0.0, 0.0, 1.0, //bottom left vertex, top left texture
		1.0, 1.0, 1.0, 0.0, //top right, bottom right texture
	}
	//Indices defining which vertices to use (useful to avoid specifying duplicate vertices in the array)
	indices := []int32{
		0, 1, 2,
		0, 3, 1,
	}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	//Work on this specific object
	gl.BindVertexArray(vao)

	//Load the vertices into memory
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	//Load the indices into memory
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	//Position and texture attributes (the "location = 0" in the shader)
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
	//Adapter between the engine's listener and GLFW's listener
	w.glfwWindow.SetKeyCallback(func(w *glfw.Window, glfwKey glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		k := translateKey(glfwKey)
		if k == key.UNKNOWN {
			return //Ignore unknown keys
		}

		if action == glfw.Press {
			listener.KeyPressed(k)
		} else if action == glfw.Release {
			listener.KeyReleased(k)
		}
	})
}

func (w *windowImpl) SetCursorListener(cursorListener cursor.Listener) {
	//Adapter between the engine's listener and GLFW's listener
	w.glfwWindow.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		cursorListener.CursorMoved(int(xpos), int(ypos))
	})
}

func (w *windowImpl) CreateImageLoader() ui.ImageLoader {
	return &imageLoaderImpl{}
}

func (w *windowImpl) CreateGraphics() ui.Graphics {
	//Clear the screen before delegating the drawing to the current state
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	return &graphicsImpl{shaderProgram: w.shaderProgram}
}

func (w *windowImpl) ShouldClose() bool {
	return w.glfwWindow.ShouldClose()
}

func (w *windowImpl) Update() {
	w.glfwWindow.SwapBuffers()
	glfw.PollEvents()
}

func (w *windowImpl) Destroy() {
	//Delete the objects allocated in memory
	gl.DeleteVertexArrays(1, &w.vertexArrayObject)
	gl.DeleteBuffers(1, &w.vertexBufferObject)
	gl.DeleteBuffers(1, &w.elementBufferObject)
	gl.DeleteProgram(w.shaderProgram)

	//Release the rest of the memory
	w.glfwWindow.Destroy()
	glfw.Terminate()
}
