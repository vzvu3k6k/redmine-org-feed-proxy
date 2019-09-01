package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

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
	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")

	u, err := resolveURL(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	body, err := fetch(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

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
