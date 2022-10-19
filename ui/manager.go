package ui

import (
	"github.com/Hikarikun92/go-game-engine/key"
)

type WindowManager interface {
	CreateMainWindow() Window
}

type Window interface {
	SetKeyListener(keyListener key.Listener)

	ShouldClose() bool
	Update()
	Destroy()
}
