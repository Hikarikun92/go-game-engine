package gl

import (
	"fmt"
	"github.com/Hikarikun92/go-game-engine/key"
	"github.com/Hikarikun92/go-game-engine/ui"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/draw"
	_ "image/jpeg"
	"log"
	"os"
	"runtime"
	"strings"
)

type glWindowManager struct {
}

func NewWindowManager() ui.WindowManager {
	return &glWindowManager{}
}

type windowImpl struct {
	glfwWindow         *glfw.Window
	vertexArrayObject  uint32
	vertexBufferObject uint32
	shaderProgram      uint32
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
		-0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 1.0, 1.0,
	}

	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	//position and texture attributes (the "location = 0" in the shader)
	gl.VertexAttribPointerWithOffset(0, 4, gl.FLOAT, false, 4*4 /* 4 values per vertex * 4 bytes per value */, 0)
	gl.EnableVertexAttribArray(0)

	window.Show()

	return &windowImpl{
		glfwWindow:         window,
		vertexArrayObject:  vao,
		vertexBufferObject: vbo,
		shaderProgram:      shaderProgram,
	}
}

var vertexShader = `
#version 330 core
layout (location = 0) in vec4 vertexData;

out vec2 TexCoord;

void main()
{
	gl_Position = vec4(vertexData.xy, 0.0, 1.0);
	TexCoord = vertexData.zw;
}
` + "\x00"

var fragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 TexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, TexCoord);
}
` + "\x00"

func newShaderProgram() (uint32, error) {
	vertexShader, err := compileShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logValue := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logValue))

		return 0, fmt.Errorf("failed to link program: %v", logValue)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	cSources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cSources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logValue := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logValue))

		return 0, fmt.Errorf("failed to compile %v: %v", source, logValue)
	}

	return shader, nil
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

func translateKey(glfwKey glfw.Key) key.Key {
	switch glfwKey {
	case glfw.KeySpace:
		return key.SPACE
	case glfw.KeyApostrophe:
		return key.APOSTROPHE
	case glfw.KeyComma:
		return key.COMMA
	case glfw.KeyMinus:
		return key.MINUS
	case glfw.KeyPeriod:
		return key.PERIOD
	case glfw.KeySlash:
		return key.SLASH
	case glfw.Key0:
		return key.ZERO
	case glfw.Key1:
		return key.ONE
	case glfw.Key2:
		return key.TWO
	case glfw.Key3:
		return key.THREE
	case glfw.Key4:
		return key.FOUR
	case glfw.Key5:
		return key.FIVE
	case glfw.Key6:
		return key.SIX
	case glfw.Key7:
		return key.SEVEN
	case glfw.Key8:
		return key.EIGHT
	case glfw.Key9:
		return key.NINE
	case glfw.KeySemicolon:
		return key.SEMICOLON
	case glfw.KeyEqual:
		return key.EQUAL
	case glfw.KeyA:
		return key.A
	case glfw.KeyB:
		return key.B
	case glfw.KeyC:
		return key.C
	case glfw.KeyD:
		return key.D
	case glfw.KeyE:
		return key.E
	case glfw.KeyF:
		return key.F
	case glfw.KeyG:
		return key.G
	case glfw.KeyH:
		return key.H
	case glfw.KeyI:
		return key.I
	case glfw.KeyJ:
		return key.J
	case glfw.KeyK:
		return key.K
	case glfw.KeyL:
		return key.L
	case glfw.KeyM:
		return key.M
	case glfw.KeyN:
		return key.N
	case glfw.KeyO:
		return key.O
	case glfw.KeyP:
		return key.P
	case glfw.KeyQ:
		return key.Q
	case glfw.KeyR:
		return key.R
	case glfw.KeyS:
		return key.S
	case glfw.KeyT:
		return key.T
	case glfw.KeyU:
		return key.U
	case glfw.KeyV:
		return key.V
	case glfw.KeyW:
		return key.W
	case glfw.KeyX:
		return key.X
	case glfw.KeyY:
		return key.Y
	case glfw.KeyZ:
		return key.Z
	case glfw.KeyLeftBracket:
		return key.LEFT_BRACKET
	case glfw.KeyBackslash:
		return key.BACKSLASH
	case glfw.KeyRightBracket:
		return key.RIGHT_BRACKET
	case glfw.KeyGraveAccent:
		return key.GRAVE_ACCENT
	case glfw.KeyWorld1:
		return key.WORLD_1
	case glfw.KeyWorld2:
		return key.WORLD_2
	case glfw.KeyEscape:
		return key.ESCAPE
	case glfw.KeyEnter:
		return key.ENTER
	case glfw.KeyTab:
		return key.TAB
	case glfw.KeyBackspace:
		return key.BACKSPACE
	case glfw.KeyInsert:
		return key.INSERT
	case glfw.KeyDelete:
		return key.DELETE
	case glfw.KeyRight:
		return key.RIGHT
	case glfw.KeyLeft:
		return key.LEFT
	case glfw.KeyDown:
		return key.DOWN
	case glfw.KeyUp:
		return key.UP
	case glfw.KeyPageUp:
		return key.PAGE_UP
	case glfw.KeyPageDown:
		return key.PAGE_DOWN
	case glfw.KeyHome:
		return key.HOME
	case glfw.KeyEnd:
		return key.END
	case glfw.KeyCapsLock:
		return key.CAPS_LOCK
	case glfw.KeyScrollLock:
		return key.SCROLL_LOCK
	case glfw.KeyNumLock:
		return key.NUM_LOCK
	case glfw.KeyPrintScreen:
		return key.PRINT_SCREEN
	case glfw.KeyPause:
		return key.PAUSE
	case glfw.KeyF1:
		return key.F1
	case glfw.KeyF2:
		return key.F2
	case glfw.KeyF3:
		return key.F3
	case glfw.KeyF4:
		return key.F4
	case glfw.KeyF5:
		return key.F5
	case glfw.KeyF6:
		return key.F6
	case glfw.KeyF7:
		return key.F7
	case glfw.KeyF8:
		return key.F8
	case glfw.KeyF9:
		return key.F9
	case glfw.KeyF10:
		return key.F10
	case glfw.KeyF11:
		return key.F11
	case glfw.KeyF12:
		return key.F12
	case glfw.KeyF13:
		return key.F13
	case glfw.KeyF14:
		return key.F14
	case glfw.KeyF15:
		return key.F15
	case glfw.KeyF16:
		return key.F16
	case glfw.KeyF17:
		return key.F17
	case glfw.KeyF18:
		return key.F18
	case glfw.KeyF19:
		return key.F19
	case glfw.KeyF20:
		return key.F20
	case glfw.KeyF21:
		return key.F21
	case glfw.KeyF22:
		return key.F22
	case glfw.KeyF23:
		return key.F23
	case glfw.KeyF24:
		return key.F24
	case glfw.KeyF25:
		return key.F25
	case glfw.KeyKP0:
		return key.KEYPAD_0
	case glfw.KeyKP1:
		return key.KEYPAD_1
	case glfw.KeyKP2:
		return key.KEYPAD_2
	case glfw.KeyKP3:
		return key.KEYPAD_3
	case glfw.KeyKP4:
		return key.KEYPAD_4
	case glfw.KeyKP5:
		return key.KEYPAD_5
	case glfw.KeyKP6:
		return key.KEYPAD_6
	case glfw.KeyKP7:
		return key.KEYPAD_7
	case glfw.KeyKP8:
		return key.KEYPAD_8
	case glfw.KeyKP9:
		return key.KEYPAD_9
	case glfw.KeyKPDecimal:
		return key.KEYPAD_DECIMAL
	case glfw.KeyKPDivide:
		return key.KEYPAD_DIVIDE
	case glfw.KeyKPMultiply:
		return key.KEYPAD_MULTIPLY
	case glfw.KeyKPSubtract:
		return key.KEYPAD_SUBTRACT
	case glfw.KeyKPAdd:
		return key.KEYPAD_ADD
	case glfw.KeyKPEnter:
		return key.KEYPAD_ENTER
	case glfw.KeyKPEqual:
		return key.KEYPAD_EQUAL
	case glfw.KeyLeftShift:
		return key.LEFT_SHIFT
	case glfw.KeyLeftControl:
		return key.LEFT_CONTROL
	case glfw.KeyLeftAlt:
		return key.LEFT_ALT
	case glfw.KeyLeftSuper:
		return key.LEFT_SUPER
	case glfw.KeyRightShift:
		return key.RIGHT_SHIFT
	case glfw.KeyRightControl:
		return key.RIGHT_CONTROL
	case glfw.KeyRightAlt:
		return key.RIGHT_ALT
	case glfw.KeyRightSuper:
		return key.RIGHT_SUPER
	case glfw.KeyMenu:
		return key.MENU
	default:
		return key.UNKNOWN
	}
}

func (w *windowImpl) CreateImageLoader() ui.ImageLoader {
	return &imageLoaderImpl{}
}

type imageLoaderImpl struct {
}

type imageImpl struct {
	textureId uint32
}

func (i *imageLoaderImpl) LoadImage(file string) ui.Image {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v", file, err)
	}
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

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgbaSize.X), int32(rgbaSize.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return imageImpl{textureId: texture}
}

func (i *imageLoaderImpl) UnloadImage(image ui.Image) {
	img := image.(imageImpl)
	gl.DeleteTextures(1, &img.textureId)
}

func (w *windowImpl) CreateGraphics() ui.Graphics {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	return &graphicsImpl{}
}

type graphicsImpl struct {
}

func (g *graphicsImpl) DrawImage(image ui.Image, x int, y int) {
	img := image.(imageImpl)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, img.textureId)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
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
	gl.DeleteProgram(w.shaderProgram)

	w.glfwWindow.Destroy()
	glfw.Terminate()
}
