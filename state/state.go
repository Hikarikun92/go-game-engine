package state

import (
	"github.com/Hikarikun92/go-game-engine/ui"
	"time"
)

type State interface {
	Load(imageLoader ui.ImageLoader)
	Update(delta time.Duration) State
	Draw(graphics ui.Graphics)
	Unload(imageLoader ui.ImageLoader)
}
