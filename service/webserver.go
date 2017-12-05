package service

import (
	"log"
	"net/http"
	"time"
)

// StartWebServer start HTTP server listening at the given port
func StartWebServer(port string) {
	log.Println("Starting HTTP server at " + port)
	r := NewRouter()
	// http.Handle("/", r)
	svr := http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	svr.ListenAndServe()
	err := svr.ListenAndServe()

	if err != nil {
		log.Println("An error occured starting HTTP listner at port " + port)
		log.Println(err.Error())
	}

}
