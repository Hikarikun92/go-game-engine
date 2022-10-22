package ui

import (
	"github.com/Hikarikun92/go-game-engine/cursor"
	"github.com/Hikarikun92/go-game-engine/key"
	"github.com/Hikarikun92/go-game-engine/settings"
)

type WindowManager interface {
	CreateMainWindow(settings *settings.Settings) Window
}

type Window interface {
	SetKeyListener(keyListener key.Listener)
	SetCursorListener(cursorListener cursor.Listener)

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
