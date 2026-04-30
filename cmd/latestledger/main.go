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
	log.Fatal(http.ListenAndServe(addr, mux))
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
