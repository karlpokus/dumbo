package dumbo

import (
	"bytes"
	"io"
	"sync"

	"dumbo/internal/file"
)

type Data struct {
	gz    []byte
	hash  string
	store Store
	sync.Mutex
}

// Save splits the compressed stream into an in-memory copy and decompressed to disk
func (data *Data) Save(r io.Reader) error {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	data.store.Reset()
	if err := Decompress(data.store, tee); err != nil {
		return err
	}
	if v, ok := data.store.(*file.Store); ok {
		if err := v.Sync(); err != nil {
			return err
		}
	}
	b := buf.Bytes()
	data.gz = b
	data.hash = hash(b)
	return nil
}

func FileStore(fpath string) (*file.Store, error) {
	return file.New(fpath)
}

// New returns a ready-to-use Data type
func New(store Store) (*Data, error) {
	var buf bytes.Buffer
	err := compress(&buf, store)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()
	return &Data{
		gz:    b,
		store: store,
		hash:  hash(b),
	}, nil
}
