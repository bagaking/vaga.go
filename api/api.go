package api

import (
	"github.com/bagaking/vaga.go/localVideos"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers(config []localVideos.VideoBlobConf) *httprouter.Router {
	blobs := make([]*localVideos.VideoBlob, len(config))
	for i, conf := range config {
		blobs[i] = &localVideos.VideoBlob{VideoBlobConf: conf}
	}
	initial(blobs)

	router := httprouter.New()

	router.GET("/", HandlerVideoIndex)
	router.GET("/tree/:blob_ind", HandlerVideoTree)
	router.GET("/watch/:video_hash", HandlerVideoWatch)
	router.GET("/video/:video_hash", HandlerVideoStream)
	return router
}
