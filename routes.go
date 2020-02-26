package dumbo

import (
	"log"
	"net/http"
)

func Read(data *Data) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		data.Lock()
		defer data.Unlock()
		data.Send(w)
	}
}

func Write(data *Data, stderr *log.Logger) http.HandlerFunc {
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
			stderr.Printf("Unable to save request body: %s\n", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.WriteHeader(201)
	}
}
