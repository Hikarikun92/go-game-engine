package ui

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
)

type Window interface {
	ShouldClose() bool
	Update()
	Destroy()
}

type windowImpl struct {
	glfwWindow *glfw.Window
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func CreateMainWindow() Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(800, 600, "Example game", nil, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}
	window.MakeContextCurrent()

	return &windowImpl{glfwWindow: window}
}

func (w *windowImpl) ShouldClose() bool {
	return w.glfwWindow.ShouldClose()
}

func (w *windowImpl) Update() {
	w.glfwWindow.SwapBuffers()
	glfw.PollEvents()
}

func (w *windowImpl) Destroy() {
	w.glfwWindow.Destroy()
	glfw.Terminate()
}
