package dumbo

import (
	"log"
	"net/http"

	"github.com/karlpokus/ratelmt"
)

func Routes(data *Data, stderr *log.Logger, limit int) http.Handler {
	router := http.NewServeMux()
	router.Handle("/read", ratelmt.Mw(float64(limit), read(data)))
	router.Handle("/write", ratelmt.Mw(float64(limit), write(data, stderr)))
	return router
}

func read(data *Data) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		if r.Header.Get("Etag") == data.hash {
			w.WriteHeader(304)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		data.Lock()
		defer data.Unlock()
		w.Header().Set("Etag", data.hash)
		w.Write(data.gz)
	}
}

func write(data *Data, stderr *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		if r.Header.Get("Content-Encoding") != "gzip" {
			http.Error(w, "Invalid or missing encoding header", 400)
			return
		}
		if r.ContentLength <= 0 {
			http.Error(w, "Empty or unknown content length", 400)
			return
		}
		defer r.Body.Close()
		data.Lock() // ok to lock twice. blocks until available.
		defer data.Unlock()
		err := data.Save(r.Body)
		if err != nil {
			if stderr != nil {
				stderr.Printf("Unable to save request body: %s\n", err)
			}
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Etag", data.hash)
		w.WriteHeader(201)
	}
}
