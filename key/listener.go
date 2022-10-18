package key

type Listener interface {
	KeyPressed(key Key)
	KeyReleased(key Key)
}
