package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/time/rate"
)

// allows only one request per 50 seconds (allowing some bursts)
var limiter = rate.NewLimiter(0.02, 3)

func resolveURL(path string) (*url.URL, error) {
	u, err := url.Parse("https://redmine.org/" + path)
	if err != nil {
		return nil, errors.New("invalid path")
	}

	if u.Hostname() != "redmine.org" {
		return nil, errors.New("hostname must be redmine.org")
	}

	if !strings.HasSuffix(u.RequestURI(), ".atom") {
		return nil, errors.New("path must be ended with `.atom`")
	}

	return u, nil
}

func fetch(u *url.URL) (io.ReadCloser, error) {
	resp, err := http.Get(u.String())

	// TODO: return a proper status code according to resp.StatusCode
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("can't fetch a requested URL")
	}
	return resp.Body, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	u, err := resolveURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if limiter.Allow() == false {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	body, err := fetch(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
	io.Copy(w, body)
	body.Close()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
