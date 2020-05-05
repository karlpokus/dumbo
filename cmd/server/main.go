package main

import (
	"log"
	"os"

	"dumbo"
	"github.com/karlpokus/srv"
)

func main() {
  stdout := log.New(os.Stdout, "server ", log.Ldate|log.Ltime)
  stderr := log.New(os.Stderr, "server ", log.Ldate|log.Ltime)
  if len(os.Args) == 1 {
    stderr.Fatal("Missing fpath arg")
  }
  fpath := os.Args[1]
	s, err := srv.New(func(s *srv.Server) error {
    f, err := dumbo.FileStore(fpath)
    if err != nil {
			return err
		}
		data, err := dumbo.New(f)
		if err != nil {
			return err
		}
		s.Router = dumbo.Routes(data, stderr, 1)
		s.Logger = stdout
		s.Host = "0.0.0.0"
		s.Port = "7979"
		return nil
	})
	if err != nil {
		stderr.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		stderr.Fatal(err)
	}
	stdout.Println("main exited")
}
