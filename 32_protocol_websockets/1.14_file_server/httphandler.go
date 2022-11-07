package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func main() {
	port := ":8083"
	http.Handle("/count", new(countHandler))
	fmt.Printf("port start:%v\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
