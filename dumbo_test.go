package dumbo

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/karlpokus/routest/v2"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%s", err)
	}
}

// a should match b
func TestPersistance(t *testing.T) {
	a := "this is nice"
	store := bytes.NewBuffer([]byte(a))
	data, err := New(store)
	fatal(t, err)
	err = data.Save(bytes.NewBuffer(data.gz))
	fatal(t, err)
	b := store.String()
	if a != b {
		fatal(t, fmt.Errorf("Expected %s and %s to be equal", a, b))
	}
}

func TestRoutes(t *testing.T) {
	str := "important data"
	store := bytes.NewBuffer([]byte(str))
	data, err := New(store)
	fatal(t, err)

	update := bytes.NewBuffer([]byte("even more important data"))
	var updateComp bytes.Buffer
	err = compress(&updateComp, update)
	fatal(t, err)

	routest.Test(t, func() http.Handler {
		return Routes(data, nil, 10)
	}, []routest.Data{
		{
			Name:    "read method not allowed",
			Method:  "PUT",
			Path:    "/read",
			Status:  405,
		},
		{
			Name:         "read stateless",
			Method:       "GET",
			Path:         "/read",
			Status:       200,
			ResponseBody: data.gz,
			ResponseHeader: http.Header{
				"Content-Encoding": []string{"gzip"},
				"Etag": []string{data.hash},
			},
		},
		{
			Name:   "read etag match",
			Method: "GET",
			Path:   "/read",
			RequestHeader: http.Header{
				"Etag": []string{data.hash},
			},
			Status:  304,
			ResponseBody: nil,
		},
		{
			Name:    "write method not allowed",
			Method:  "GET",
			Path:    "/write",
			Status:  405,
		},
		{
			Name:    "write missing header",
			Method:  "POST",
			Path:    "/write",
			Status:  400,
		},
		{
			Name:    "write zero content length",
			Method:  "POST",
			Path:    "/write",
			RequestHeader: http.Header{
				"Content-Encoding": []string{"gzip"},
			},
			Status:  400,
		},
		{
			Name: "write valid request body",
			Method: "POST",
			Path: "/write",
			RequestHeader: http.Header{
				"Content-Encoding": []string{"gzip"},
			},
			RequestBody: &updateComp,
			Status: 201,
		},
	})
}
