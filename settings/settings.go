package settings

type Settings struct {
	Width       int
	Height      int
	WindowTitle string
	Fps         int
}

func DefaultSettings() *Settings {
	return &Settings{
		Width:       800,
		Height:      600,
		WindowTitle: "Example game",
		Fps:         60,
	}
}
