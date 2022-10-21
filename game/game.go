package game

import (
	"github.com/Hikarikun92/go-game-engine/key"
	"github.com/Hikarikun92/go-game-engine/state"
	"github.com/Hikarikun92/go-game-engine/ui"
	"log"
	"time"
)

type Game interface {
	Start()
}

type gameImpl struct {
	windowManager ui.WindowManager
	state         state.State
}

func NewGame(windowManager ui.WindowManager, initialState state.State) Game {
	return &gameImpl{windowManager: windowManager, state: initialState}
}

func (game *gameImpl) Start() {
	window := game.windowManager.CreateMainWindow()
	defer window.Destroy()

	window.SetKeyListener(game)

	imageLoader := window.CreateImageLoader()

	running := true
	game.state.Load(imageLoader)

	previousTime := time.Now()
	ticker := time.NewTicker(16667 * time.Microsecond)

	for running {
		if window.ShouldClose() {
			ticker.Stop()
			game.state.Unload(imageLoader)
			running = false
			break
		}

		select {
		case t := <-ticker.C:
			delta := t.Sub(previousTime)

			nextState := game.state.Update(delta)

			graphics := window.CreateGraphics()
			game.state.Draw(graphics)

			if nextState == nil {
				ticker.Stop()
				game.state.Unload(imageLoader)
				running = false
			} else if nextState != game.state {
				game.state.Unload(imageLoader)

				game.state = nextState
				game.state.Load(imageLoader)
			}

			previousTime = t
		}

		window.Update()
	}
}

func (game *gameImpl) KeyPressed(k key.Key) {
	log.Println("KeyPressed:", k)

	listener, isListener := game.state.(key.Listener)
	if isListener {
		listener.KeyPressed(k)
	}
}

func (game *gameImpl) KeyReleased(k key.Key) {
	log.Println("KeyReleased:", k)

	listener, isListener := game.state.(key.Listener)
	if isListener {
		listener.KeyReleased(k)
	}
}
