package service

import (
	"fmt"
	"net/http"
	"time"
)

// StartWebServer start HTTP server listening at the given port
func StartWebServer(port string) {
	fmt.Println("Starting HTTP server at " + port)
	r := NewRouter()
	// http.Handle("/", r)
	svr := http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	err := svr.ListenAndServe()
	if err != nil {
		fmt.Printf("HTTP server failed, port %s, err: %s\n", port, err.Error())
	}

}
