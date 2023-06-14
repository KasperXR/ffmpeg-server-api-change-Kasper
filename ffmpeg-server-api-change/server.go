package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

const IO_DIR = "/usr/local/etc/"
const ADDR = ":443"
const DEBUGADDR = ":8080"

func StartServer(debug bool) {
	mux := http.NewServeMux()

	startFileServer(mux)

	handleAPICall(mux)

	srv := makeConfigs(mux, debug)

	// Start the server
	if debug {
		fmt.Println("\nDEBUG MODE\n\n--------\n\nPORT 80 ONLY\n\n--------")
		log.Fatal(srv.ListenAndServe())
	} else {
		log.Fatal(srv.ListenAndServeTLS(
			"/etc/letsencrypt/live/api.mit-hjerte.dk/fullchain.pem",
			"/etc/letsencrypt/live/api.mit-hjerte.dk/privkey.pem"))
	}
}

func startFileServer(mux *http.ServeMux) {
	videos := http.FileServer(http.Dir("./videos"))

	mux.Handle("/videos/", http.StripPrefix("/videos/", addHeaders(videos)))
}

func makeConfigs(mux *http.ServeMux, debug bool) *http.Server {
	if !debug {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		//Config server
		return &http.Server{
			Addr:         ADDR,
			Handler:      mux,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
	} else {
		return &http.Server{
			Addr:    DEBUGADDR,
			Handler: mux,
		}
	}
}

func handleAPICall(mux *http.ServeMux) {
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// Read all the headers of the request and log them
		// m := readAllHeaders(r)

		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Add("Cache-Control", "no-cache, must-revalidate, proxy-revalidate")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		var requestJSON JSONObj
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&requestJSON)

		uuid := uuid.New()

		if err != nil {
			fmt.Println("JSON Decode error:", err)
			return
		} else {
			fmt.Println("JSON Decode success")

			// Get current date in format DD-MM-YYYY
			date := time.Now().Format("02-01-2006")

			GenerateVideo("mit-hjerte-"+date+"-"+uuid.String(), requestJSON.Payload)

			// Send a URL for downloading the file, to the client
			response := "https://mit-hjerte.dk/download?url=" + "mit-hjerte-" + date + "-" + uuid.String() + "-final.mp4"

			// Send video URL as response for use in front-end
			w.Write([]byte(response))
		}

	})
}

func HandleIndex(mux *http.ServeMux) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// This `http.ResponseWriter` and `*http.Request` are standard
		// library objects that are used to respond to HTTP requests.

		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		// Read all the headers of the request and log them
		// The `readAllHeaders` function returns the headers as a map.
		m := readAllHeaders(r)

		// The `fmt.Println` function prints a string to the console.
		fmt.Println("All Headers:", m)

		serveFile(w, r, "index.html")
	})

}

// serveFile is a helper function that sends a file to the client.
func serveFile(w http.ResponseWriter, r *http.Request, f string) {
	// The `http.ServeFile` function sends a file to the client.
	// The first argument is the `http.ResponseWriter` that will send
	// the file to the client.
	// The second argument is the `http.Request` that contains the
	// information about the request.
	// The third argument is the path to the file.
	http.ServeFile(w, r, f)
}

// readHeader is a helper function that reads the header from a request
// and returns the value of the header.
func readHeader(r *http.Request, h string) string {
	// The `r.Header` function returns the header of the request.
	// The `http.Header` type is a map of strings to strings.
	// The `http.Header.Get` function returns the value of the
	// header with the given name.
	return r.Header.Get(h)
}

// readAllHeaders is a helper function that reads all the headers from a
// request and returns them as a map.
func readAllHeaders(r *http.Request) map[string]string {
	// The `r.Header` function returns the header of the request.
	// The `http.Header` type is a map of strings to strings.
	// The `http.Header.Get` function returns the value of the
	// header with the given name.
	h := r.Header
	// The `make` function returns a map with the given size.
	m := make(map[string]string)
	// The loop iterates over the map and
	// calls the given function for each key-value pair.
	for k, v := range h {
		// The `m[k]` function returns the value of the key.
		m[k] = v[0]
	}
	return m
}

func addHeaders(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Add("Cache-Control", "no-cache, must-revalidate, proxy-revalidate")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	}
}

// readBody is a helper function that reads the body of a request
// and returns the value of the body as a json string.
func readBody(r *http.Request) string {
	// The `r.Body` function returns the body of the request.
	// The `io.ReadAll` function reads all the bytes from the body
	// and returns them as a string.
	b, _ := io.ReadAll(r.Body)

	return string(b)
}

// This function returns true if the request body has a key with the given value
// This can be expanded to accept standard form data as that's what the client sends
// Currently it's just a placeholder copilot helped make to test the server's handling of files
func bodyHasKey(r *http.Request, key string, value string) bool {
	// The `r.Body` function returns the body of the request.
	// The `io.ReadAll` function reads all the bytes from the body
	// and returns them as a string.
	b, _ := io.ReadAll(r.Body)

	// The `json.Unmarshal` function decodes a JSON string into
	// a value.
	// The first argument is the JSON string.
	// The second argument is the value to decode the JSON string into.
	var m map[string]interface{}
	json.Unmarshal(b, &m)

	// The `m[key]` function returns the value of the key.
	return m[key] == value
}
