package state

import (
	"github.com/Hikarikun92/go-game-engine/graphics"
	"time"
)

type State interface {
	Load()
	Update(delta time.Duration) State
	Draw(graphics graphics.Graphics)
	Unload()
}
