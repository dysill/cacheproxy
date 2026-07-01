package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type item struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	resp := item{
		ID:          id,
		Description: "Fast data for id: " + id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func slowHandler(delay time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// Simulate data retrieval delay
		time.Sleep(delay)
		resp := item{
			ID:          id,
			Description: "Slow data for id: " + id,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	port := flag.Int("port", 8080, "Port number for upstream")
	delayInt := flag.Int("delay", 200, "Delay for slow data retrieval in ms")
	flag.Parse()
	delay := time.Duration(*delayInt) * time.Millisecond

	mux := http.NewServeMux()
	mux.HandleFunc("/fast/{id}", fastHandler)
	mux.HandleFunc("/slow/{id}", slowHandler(delay))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("upstream server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))

}
