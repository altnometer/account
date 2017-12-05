package main

import (
	"flag"

	"github.com/altnometer/account/service"
)

func main() {
	port := flag.String("port", "8080", "server port")
	flag.Parse()
	service.StartWebServer(*port)
}
