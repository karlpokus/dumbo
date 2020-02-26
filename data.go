package dumbo

import (
	"bytes"
	"io"
	"sync"

	"dumbo/internal/file"
)

type Data struct {
	gz    []byte
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
	data.gz = buf.Bytes()
	return nil
}

// Send writes the compressed blob to w
func (data *Data) Send(w io.Writer) {
	w.Write(data.gz)
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
	return &Data{
		gz:    buf.Bytes(),
		store: store,
	}, nil
}
