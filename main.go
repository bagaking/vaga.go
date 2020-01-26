package main

import (
	"github.com/bagaking/vaga.go/api"
	"net/http"
)

const (
	Port = ":9001"
)

func main() {
	r := api.RegisterHandlers(AllAvailableVideoBlobs)
	http.ListenAndServe(Port, r)
}
