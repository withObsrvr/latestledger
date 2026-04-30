package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"latestledger.com/latestledger/internal/app"
	"latestledger.com/latestledger/internal/templates"
)

func main() {
	client := app.NewClient(os.Getenv("OBSRVR_API_KEY"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data := loadPageData(r, client)
		if err := templates.Layout(data).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		data := loadPageData(r, client)
		if err := templates.App(data).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/fragment", func(w http.ResponseWriter, r *http.Request) {
		data := loadPageData(r, client)
		if err := templates.Dashboard(data).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("latestledger listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, requestLogger(noCache(mux))))
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		if lrw.status == 0 {
			lrw.status = http.StatusOK
		}
		log.Printf("%s %s %s %d %dB %s remote=%s ua=%q", r.Method, r.URL.RequestURI(), r.Proto, lrw.status, lrw.bytes, time.Since(start).Round(time.Millisecond), r.RemoteAddr, r.UserAgent())
	})
}

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

func loadPageData(r *http.Request, client *app.Client) app.PageData {
	network := app.ParseNetwork(r.URL.Query().Get("network"))
	data := app.PageData{Network: network, Now: time.Now().UTC()}
	stats, err := client.NetworkStats(r.Context(), network)
	if err != nil {
		data.Error = err.Error()
		return data
	}
	data.Stats = stats
	return data
}
