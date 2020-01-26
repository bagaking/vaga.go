package api

import (
	"github.com/bagaking/vaga.go/localVideos"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(allAvailableVideoBlobs []*localVideos.VideoBlob) *httprouter.Router {
	initial(allAvailableVideoBlobs)

	router := httprouter.New()

	router.GET("/", HandlerVideoIndex)
	router.GET("/tree/:blob_ind", HandlerVideoTree)
	router.GET("/watch/:video_hash", HandlerVideoWatch)
	router.GET("/video/:video_hash", HandlerVideoStream)
	return router
}



