package game

import (
	"github.com/Hikarikun92/go-game-engine/state"
	"github.com/Hikarikun92/go-game-engine/ui"
	"time"
)

func Start(currentState state.State) {
	window := ui.CreateMainWindow()
	defer window.Destroy()

	running := true
	currentState.Load()

	previousTime := time.Now()
	ticker := time.NewTicker(16667 * time.Microsecond)

	for running {
		if window.ShouldClose() {
			ticker.Stop()
			currentState.Unload()
			running = false
			break
		}

		select {
		case t := <-ticker.C:
			delta := t.Sub(previousTime)

			nextState := currentState.Update(delta)

			//currentState.Draw(graphics)

			if nextState == nil {
				ticker.Stop()
				currentState.Unload()
				running = false
			} else if nextState != currentState {
				currentState.Unload()

				currentState = nextState
				currentState.Load()
			}

			previousTime = t
		}

		window.Update()
	}
}
