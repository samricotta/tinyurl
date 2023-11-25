package server

import (
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
		case http.MethodGet: // localhost:8080/gz (2 byte request)
			// actualURL := getActualURL(r.URL.Path) // path is gz (the 2 byte ID)
			// // we get the actual URL from the database and redirect the user
			store.Get(r.URL.Path)
			// http.Redirect(w, r, actualURL, http.StatusFound)

		case http.MethodPost:
			//body is currently in bytes
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			tinyUrl, err := store.Set(bodyBytes)
			if err != nil {
				panic(err)
			}
			// we return the tiny URL to the user

			w.Write(tinyUrl)
			
		default:
		}
	})

	http.ListenAndServe("localhost:8080", nil)
}
