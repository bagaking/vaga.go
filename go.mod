module github.com/bagaking/vaga.go

go 1.14

replace github.com/bagaking/vaga.go/localVideos => ./localVideos

replace github.com/bagaking/vaga.go/api => ./api

replace github.com/bagaking/vaga.go/conf => ./conf

require (
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/json-iterator/go v1.1.10
	github.com/julienschmidt/httprouter v1.3.0
)
