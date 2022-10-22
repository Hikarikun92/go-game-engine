package cursor

type Listener interface {
	CursorMoved(x int, y int)
}
