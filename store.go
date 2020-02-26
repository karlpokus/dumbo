package dumbo

type Store interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Reset()
}
