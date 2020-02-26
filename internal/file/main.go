package file

import (
	"os"
)

type Store struct {
	file *os.File
}

func (st *Store) Read(p []byte) (n int, err error) {
	return st.file.Read(p)
}

func (st *Store) Write(p []byte) (n int, err error) {
	return st.file.Write(p)
}

func (st *Store) Reset() { // would like to return an error here
	st.file.Seek(0, 0)
	st.file.Truncate(0)
}

// Sync flushes bits to disk
func (st *Store) Sync() error {
	return st.file.Sync()
}

func New(fpath string) (*Store, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return &Store{
		file: f,
	}, nil
}
