package main

import (
	"bagaking.com/vaga.go/api"
	"net/http"
)

const (
	Port = ":9001"
)

func main() {
	r := api.RegisterHandlers(AllAvailableVideoBlobs)
	http.ListenAndServe(Port, r)
}
