package game

import (
	"github.com/Hikarikun92/go-game-engine/key"
	"github.com/Hikarikun92/go-game-engine/state"
	"github.com/Hikarikun92/go-game-engine/ui"
	"log"
	"time"
)

type Game struct {
	windowManager ui.WindowManager
	state         state.State
}

func NewGame(windowManager ui.WindowManager, initialState state.State) *Game {
	return &Game{windowManager: windowManager, state: initialState}
}

func (game *Game) Start() {
	window := game.windowManager.CreateMainWindow()
	defer window.Destroy()

	window.SetKeyListener(game)

	running := true
	game.state.Load()

	previousTime := time.Now()
	ticker := time.NewTicker(16667 * time.Microsecond)

	for running {
		if window.ShouldClose() {
			ticker.Stop()
			game.state.Unload()
			running = false
			break
		}

		select {
		case t := <-ticker.C:
			delta := t.Sub(previousTime)

			nextState := game.state.Update(delta)

			//currentState.Draw(graphics)

			if nextState == nil {
				ticker.Stop()
				game.state.Unload()
				running = false
			} else if nextState != game.state {
				game.state.Unload()

				game.state = nextState
				game.state.Load()
			}

			previousTime = t
		}

		window.Update()
	}
}

func (game *Game) KeyPressed(k key.Key) {
	log.Println("KeyPressed:", k)

	listener, isListener := game.state.(key.Listener)
	if isListener {
		listener.KeyPressed(k)
	}
}

func (game *Game) KeyReleased(k key.Key) {
	log.Println("KeyReleased:", k)

	listener, isListener := game.state.(key.Listener)
	if isListener {
		listener.KeyReleased(k)
	}
}
