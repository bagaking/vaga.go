package main

import (
	"github.com/bagaking/vaga.go/api"
	"github.com/bagaking/vaga.go/conf"
	"log"
	"net/http"
)

const (
	Port = ":9001"
)

func main() {
	r := api.RegisterHandlers(conf.Static.Roots)
	log.Println("Listen at ", Port)
	http.ListenAndServe(Port, r)
}
