package game

import (
	"github.com/Hikarikun92/go-game-engine/cursor"
	"github.com/Hikarikun92/go-game-engine/key"
	"github.com/Hikarikun92/go-game-engine/settings"
	"github.com/Hikarikun92/go-game-engine/state"
	"github.com/Hikarikun92/go-game-engine/ui"
	"time"
)

type Game interface {
	Start()
}

type gameImpl struct {
	windowManager ui.WindowManager
	settings      *settings.Settings
	state         state.State
}

func NewGame(windowManager ui.WindowManager, initialState state.State, settings *settings.Settings) Game {
	return &gameImpl{windowManager: windowManager, state: initialState, settings: settings}
}

func (game *gameImpl) Start() {
	window := game.windowManager.CreateMainWindow(game.settings)
	defer window.Destroy()

	window.SetKeyListener(game)
	window.SetCursorListener(game)

	imageLoader := window.CreateImageLoader()

	running := true
	game.state.Load(imageLoader)

	previousTime := time.Now()
	ticker := time.NewTicker(1 * time.Second / time.Duration(game.settings.Fps))

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
	listener, isListener := game.state.(key.Listener)
	if isListener {
		listener.KeyPressed(k)
	}
}

func (game *gameImpl) KeyReleased(k key.Key) {
	listener, isListener := game.state.(key.Listener)
	if isListener {
		listener.KeyReleased(k)
	}
}

func (game *gameImpl) CursorMoved(x int, y int) {
	listener, isListener := game.state.(cursor.Listener)
	if isListener {
		listener.CursorMoved(x, game.settings.Height-y) //invert Y axis
	}
}
