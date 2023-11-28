package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
			path := strings.TrimPrefix(r.URL.Path, "/")
			tinyUrl, err := strconv.ParseUint(path, 10, 64)
			if err != nil {
				fmt.Println("error decoding tiny url", path, err)
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
			tinyUrlStr := strconv.FormatUint(tinyUrl, 10)

			// we return the tiny URL to the user
			path := r.URL.Host + "/" + tinyUrlStr
			fmt.Println("writing tiny url", tinyUrlStr, "path", path)
			if _, err := w.Write([]byte(path)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
		}
	})

	return http.ListenAndServe("localhost:8080", nil)
}
