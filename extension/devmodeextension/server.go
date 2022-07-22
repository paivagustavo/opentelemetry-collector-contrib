package devmode

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var (
	dbClient *dbStorageClient
	ctx      context.Context
)

// func parseSpans(rawSpans []byte) ([]Span, error) {}

func getSpansHandler(w http.ResponseWriter, r *http.Request) {
	// queryValues := r.URL.Query()
	// for when there is a query in the url, need logic to choose which param to search by
	// and check for non-empty query values
	// for now, just using GetAll

	rawSpans, err := dbClient.Get(ctx, "span_id")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	encodedSpans, err := json.Marshal(rawSpans)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	os.Stdout.Write(encodedSpans)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedSpans)
}

func startServer(ctx context.Context) error {
	var err error
	ctx = context.Background()
	dbClient, err = newClient(ctx, "sqlite3", "spans")

	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/spans", getSpansHandler)

	// setting host as always 4000 for now
	host := ":4000"
	log.Printf(`Starting server on %s`, host)
	err = http.ListenAndServe(host, mux)
	return err
}
