/*
This is a simple server that listens on two ports and is used for testing the load balancer.
This file has nothing to do with the core logic.
*/

package main

import "net/http"

func serve(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
		w.Write([]byte("Server " + port)) 
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	http.ListenAndServe(":"+port, mux)
}

func main() {
	go serve("8001")
	serve("8002")
}
