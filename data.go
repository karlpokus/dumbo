package dumbo

import (
	"bytes"
	"io"
	"os"
	"sync"
)

type Data struct {
	gz   []byte
	file *os.File
	sync.Mutex
}

// Save splits the compressed stream into an in-memory copy and decompressed to disk
func (data *Data) Save(r io.Reader) error {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	if err := reset(data.file); err != nil {
		return err
	}
	if err := Decompress(data.file, tee); err != nil {
		return err
	}
	if err := data.file.Sync(); err != nil {
		return err
	}
	data.gz = buf.Bytes()
	return nil
}

// Send writes the compressed blob to w
func (data *Data) Send(w io.Writer) {
	w.Write(data.gz)
}

// New returns a ready-to-use Data type
func New(fpath string) (*Data, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = compress(&buf, f)
	if err != nil {
		return nil, err
	}
	data := &Data{
		file: f, // keep it open
		gz:   buf.Bytes(),
	}
	return data, nil
}

// reset resets a file for writing
func reset(f *os.File) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}
	return f.Truncate(0)
}
