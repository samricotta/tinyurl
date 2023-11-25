package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/samricotta/tinyurl/store"
)

func Serve(path string) error {
	store, err := store.New(path)
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tinyUrl, err := base64.StdEncoding.DecodeString(r.URL.Path)
			if err != nil {
				fmt.Println("error decoding tiny url", r.URL.Path, err)
				http.Error(w, "Invalid URL", http.StatusBadRequest)
				return
			}
			fmt.Println("tinyUrl", tinyUrl)
			longUrl, err := store.Get(tinyUrl)
			if err != nil {
				http.Error(w, "Tiny URL not found", http.StatusNotFound)
				return
			}
			http.Redirect(w, r, string(longUrl), http.StatusFound)

		case http.MethodPost:
			//body is currently in bytes
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tinyUrl, err := store.Set(bodyBytes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// we return the tiny URL to the user
			path := r.URL.Host + "/" + base64.StdEncoding.EncodeToString(tinyUrl)
			fmt.Println("writing tiny url", tinyUrl, "path", path)
			if _, err := w.Write([]byte(path)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
		}
	})

	return http.ListenAndServe("localhost:8080", nil)
}
