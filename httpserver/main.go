package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)
func main() {
	mux := http.NewServeMux()
        mux.HandleFunc("/api/get", getValue)
        mux.HandleFunc("/api/insert", insertValue)
        mux.HandleFunc("/api/delete", deleteValue)
        mux.HandleFunc("/api/replace", insertValue)
        err := http.ListenAndServe(":3333", mux)
        if errors.Is(err, http.ErrServerClosed) {
                fmt.Printf("server closed\n")
        } else if err != nil {
                fmt.Printf("error starting server: %s\n", err)
                os.Exit(1)
        }
}
