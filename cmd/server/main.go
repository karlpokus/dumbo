package main

import (
	"log"
	"os"

	"dumbo"
	"github.com/karlpokus/ratelmt"
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
		data, err := dumbo.New(fpath)
		if err != nil {
			return err
		}
		router := s.DefaultRouter()
		router.Handle("/read", ratelmt.Mw(1, dumbo.Read(data)))
		router.Handle("/write", ratelmt.Mw(1, dumbo.Write(data, stderr)))
		s.Router = router
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
