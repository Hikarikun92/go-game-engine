package key

type Listener interface {
	KeyPressed(k Key)
	KeyReleased(k Key)
}
