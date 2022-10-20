package ui

import (
	"github.com/Hikarikun92/go-game-engine/graphics"
	"github.com/Hikarikun92/go-game-engine/key"
)

type WindowManager interface {
	CreateMainWindow() Window
}

type Window interface {
	SetKeyListener(keyListener key.Listener)

	CreateGraphics() graphics.Graphics
	ShouldClose() bool
	Update()
	Destroy()
}
