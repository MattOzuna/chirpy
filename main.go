package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	cfg := &apiConfig{
		fileserverHits: 0,
	}
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", cfg.showMetrics)
	mux.HandleFunc("GET /reset", cfg.resetMetrics)

	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serv.ListenAndServe())
}
