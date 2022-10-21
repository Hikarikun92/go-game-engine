package ui

import (
	"github.com/Hikarikun92/go-game-engine/key"
)

type WindowManager interface {
	CreateMainWindow() Window
}

type Window interface {
	SetKeyListener(keyListener key.Listener)

	CreateImageLoader() ImageLoader
	CreateGraphics() Graphics
	ShouldClose() bool
	Update()
	Destroy()
}

type ImageLoader interface {
	LoadImage(file string) Image
	UnloadImage(image Image)
}

type Image interface {
}

type Graphics interface {
	DrawImage(image Image, x int, y int)
}