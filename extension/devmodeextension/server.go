package devmodeextension

import (
	"context"
	"encoding/json"
	"errors"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"
)

func (d *devMode) getSpansHandler(w http.ResponseWriter, r *http.Request) {
	// queryValues := r.URL.Query()
	// for when there is a query in the url, need logic to choose which param to search by
	// and check for non-empty query values
	// for now, just using GetAll

	rawSpans, err := d.storage.GetAll(context.Background())
	if err != nil {
		http.NotFound(w, r)
		return
	}

	for i := range rawSpans {
		rawSpans[i].AttributesMap = make(map[string]string)
		rawSpans[i].ResourceAttributesMap = make(map[string]string)

		for _, attr := range strings.Split(rawSpans[i].Attributes, ",") {
			attrStr := strings.Split(attr, "=")
			if len(attrStr) == 2 {
				rawSpans[i].AttributesMap[attrStr[0]] = attrStr[1]
			} else {
				d.logger.Info("attribute with wrong format", zap.String("attribute", attr))
			}
		}
		rawSpans[i].Attributes = ""

		for _, attr := range strings.Split(rawSpans[i].ResourceAttributes, ",") {
			attrStr := strings.Split(attr, "=")
			if len(attrStr) == 2 {
				rawSpans[i].ResourceAttributesMap[attrStr[0]] = attrStr[1]
			} else {
				d.logger.Info("attribute with wrong format", zap.String("attribute", attr))
			}
		}
		rawSpans[i].ResourceAttributes = ""
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

func (d *devMode) startServer(ctx context.Context, logger *zap.Logger, host component.Host) error {
	var err error
	ctx = context.Background()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/spans", d.getSpansHandler)

	// setting host as always 4000 for now
	endpoint := "localhost:4000"
	log.Printf(`Starting server on %s`, endpoint)

	go func() {
		if errHTTP := http.ListenAndServe(endpoint, mux); errHTTP != nil && !errors.Is(errHTTP, http.ErrServerClosed) {
			host.ReportFatalError(errHTTP)
		}
	}()

	return err
}
